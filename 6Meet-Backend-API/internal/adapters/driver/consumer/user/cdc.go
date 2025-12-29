package user

import (
	"context"
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/constant"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/cdc"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/workerpool"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/mq/kafka"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/utils"
	"go.uber.org/zap"
)

const (
	// CDCWorkerPoolSize is the size of the worker pool for CDC events
	CDCWorkerPoolSize = 8
)

// CDCConsumer handles Debezium Change Data Capture events for user entities.
type CDCConsumer struct {
	consumer    kafka.ConsumerGroup
	userService ports.UserService
	workerPool  *workerpool.GenericPool[[]*cdc.DebeziumPayload[entity.User]]
}

// NewCDCConsumer creates a new instance of CDCConsumer.
func NewCDCConsumer(cfg *kafka.Config, userService ports.UserService) (ports.UserConsumer, error) {
	consumer, err := kafka.NewConsumer(cfg, constant.ConsumerGroupUserCDC, kafka.Recovery)
	if err != nil {
		return nil, err
	}

	// Defined fixed task function for the pool
	taskFunc := func(batch []*cdc.DebeziumPayload[entity.User]) {
		if err := userService.HandleUserBatchChange(context.Background(), batch); err != nil {
			global.Logger.Error("Failed to process batch", zap.Error(err))
		}
	}

	pool, err := workerpool.NewGenericPool(CDCWorkerPoolSize, taskFunc)
	if err != nil {
		return nil, err
	}

	return &CDCConsumer{
		consumer:    consumer,
		userService: userService,
		workerPool:  pool,
	}, nil
}

// Start starts the consumer
func (c *CDCConsumer) Start(ctx context.Context) error {
	batchSize := global.Config.Kafka.ConsumerBatchSize
	if batchSize <= 0 {
		batchSize = 100
	}

	batchInterval := utils.ToDurationMs(global.Config.Kafka.ConsumerBatchInterval)
	if batchInterval <= 0 {
		batchInterval = 500 * time.Millisecond
	}

	// Channel to buffer messages for batching.
	// Buffer size is double the batch size (batchSize * 2) to allow the handler
	// to continue pushing messages while the batch processor is flushing the
	// current batch (decoupling producer from consumer latency).
	batchChan := make(chan *cdc.DebeziumPayload[entity.User], batchSize*2)

	// Start batch processor worker
	go c.processBatchLoop(ctx, batchChan, batchSize, batchInterval)

	handler := func(ctx context.Context, key, value []byte) error {
		docID, payload, err := c.extractPayload(key, value)
		if err != nil {
			global.Logger.Warn("Skipping malformed CDC message", zap.Error(err))
			return nil
		}

		// Skip Tombstone messages (empty value for cleanup)
		if payload == nil {
			return nil
		}

		// Ensure ID is set for Delete ops (handled in extractPayload but good to be safe)
		if payload.Op == cdc.OpDelete && payload.Before != nil && payload.Before.ID == "" {
			payload.Before.ID = docID
		}

		// Push to batch channel
		select {
		case batchChan <- payload:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	}

	errHandler := func(err error) {
		global.Logger.Error("UserCDCConsumer error", zap.Error(err))
	}

	global.Logger.Info("Starting User CDC Consumer (Batch Mode)", zap.String("topic", constant.TopicUserDebezium))
	return c.consumer.Start(ctx, []string{constant.TopicUserDebezium}, handler, errHandler)
}

// processBatchLoop handles batching logic
func (c *CDCConsumer) processBatchLoop(
	ctx context.Context,
	batchChan <-chan *cdc.DebeziumPayload[entity.User],
	batchSize int,
	batchInterval time.Duration,
) {
	batch := make([]*cdc.DebeziumPayload[entity.User], 0, batchSize)
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	flush := func() {
		size := len(batch)
		if size == 0 {
			return
		}

		// Create a copy to process safely in worker pool
		finalBatch := make([]*cdc.DebeziumPayload[entity.User], size)
		copy(finalBatch, batch)

		if err := c.workerPool.Invoke(finalBatch); err != nil {
			global.Logger.Error("Failed to submit batch to worker pool", zap.Error(err))
		}

		batch = batch[:0] // Reset batch
	}

	for {
		select {
		case msg := <-batchChan:
			batch = append(batch, msg)
			if len(batch) >= batchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-ctx.Done():
			flush() // Flush remaining
			return
		}
	}
}

// extractPayload parses and converts the raw message into a domain entity payload
func (c *CDCConsumer) extractPayload(key, value []byte) (string, *cdc.DebeziumPayload[entity.User], error) {
	// Check for Tombstone message (Delete cleanup)
	if len(value) == 0 {
		return "", nil, nil
	}

	// Parse ID (Raw string as requested)
	docID := cdc.ParseMongoDBKey(key)

	// Parse into CDC DTO (handling Mongo specific types like $date)
	cdcPayload, err := cdc.ParseDebeziumMessage[CDCUser](value)
	if err != nil {
		return "", nil, err
	}

	// Map CDC DTO -> Domain Entity
	payload := &cdc.DebeziumPayload[entity.User]{
		Source: cdcPayload.Source,
		Op:     cdcPayload.Op,
		TsMs:   cdcPayload.TsMs,
	}

	if cdcPayload.After != nil {
		payload.After = cdcPayload.After.ToEntity()
	}

	// For Delete ops, extract ID from Key if "before" is missing
	if cdcPayload.Op == cdc.OpDelete {
		payload.Before = &entity.User{
			ID: docID,
		}
	} else if cdcPayload.Before != nil {
		payload.Before = cdcPayload.Before.ToEntity()
	}

	return docID, payload, nil
}

// processEvent processes a single CDC event
func (c *CDCConsumer) processEvent(ctx context.Context, evt *cdc.DebeziumPayload[entity.User]) {
	if err := c.userService.HandleUserChange(ctx, evt); err != nil {
		global.Logger.Error("Failed to handle user CDC event", zap.Error(err))
	}
}

// Stop stops the consumer
func (c *CDCConsumer) Stop() error {
	return c.consumer.Close()
}
