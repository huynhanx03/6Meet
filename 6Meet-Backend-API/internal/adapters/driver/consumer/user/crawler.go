package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/constant"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/entity"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/event"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/common/workerpool"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/mq/kafka"
	"go.uber.org/zap"
)

const (
	// WorkerPoolSize is the size of the worker pool for saving users
	WorkerPoolSize = 8
)

// CrawlerConsumer handles user crawling events from Kafka.
type CrawlerConsumer struct {
	consumer   kafka.ConsumerGroup
	repo       ports.UserRepository
	workerPool *workerpool.GenericPool[*event.UserCrawled]
}

// NewCrawlerConsumer creates a new instance of CrawlerConsumer.
func NewCrawlerConsumer(cfg *kafka.Config, repo ports.UserRepository) (ports.UserConsumer, error) {
	consumer, err := kafka.NewConsumer(cfg, constant.ConsumerGroupUserSaver, kafka.Recovery)
	if err != nil {
		return nil, err
	}

	c := &CrawlerConsumer{
		consumer: consumer,
		repo:     repo,
	}

	taskFunc := func(evt *event.UserCrawled) {
		c.processEvent(context.Background(), evt)
	}

	pool, err := workerpool.NewGenericPool(WorkerPoolSize, taskFunc)
	if err != nil {
		return nil, err
	}
	c.workerPool = pool

	return c, nil
}

// Start starts the worker to consume messages
func (c *CrawlerConsumer) Start(ctx context.Context) error {
	handler := func(ctx context.Context, key, value []byte) error {
		var evt event.UserCrawled
		if err := json.Unmarshal(value, &evt); err != nil {
			global.Logger.Error("Failed to unmarshal user crawled event", zap.Error(err))
			return nil // Skip malformed
		}

		return c.workerPool.Invoke(&evt)
	}

	errHandler := func(err error) {
		global.Logger.Error("CrawlerConsumer error", zap.Error(err))
	}

	return c.consumer.Start(ctx, []string{constant.TopicCrawlWikiPages}, handler, errHandler)
}

// processEvent processes a single user crawled event
func (c *CrawlerConsumer) processEvent(ctx context.Context, evt *event.UserCrawled) {
	saveCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &entity.User{
		Name:      evt.Name,
		Neighbors: evt.Neighbors,
		CreatedAt: evt.Timestamp,
		UpdatedAt: time.Now(),
	}

	if err := c.repo.Create(saveCtx, user); err != nil {
		global.Logger.Error("Failed to save user", zap.String("name", user.Name), zap.Error(err))
	} else {
		// global.Logger.Info("Saved user", zap.String("name", user.Name))
	}
}

// Stop stops the worker
func (c *CrawlerConsumer) Stop() error {
	return c.consumer.Close()
}
