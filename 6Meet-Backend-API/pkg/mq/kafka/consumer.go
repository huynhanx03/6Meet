package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/utils"
)

// consumerGroup wraps sarama.ConsumerGroup
type consumerGroup struct {
	client     sarama.ConsumerGroup
	topics     []string
	handler    Handler
	errHandler ErrorHandler
}

// NewConsumer creates a new Consumer Group
func NewConsumer(cfg *Config, groupID string) (ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.ClientID = cfg.ClientID

	// Rebalance Strategy
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Session.Timeout = utils.ToDurationMs(cfg.ConsumerInfo.SessionTimeout)

	client, err := sarama.NewConsumerGroup(cfg.Brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &consumerGroup{
		client: client,
	}, nil
}

// Start consumes messages from the configured topics
func (c *consumerGroup) Start(ctx context.Context, topics []string, handler Handler, errHandler ErrorHandler) error {
	c.topics = topics
	c.handler = handler
	c.errHandler = errHandler

	consumer := &consumerGroupHandler{
		handlerFunc: handler,
		errHandler:  errHandler,
	}

	// Loop to handle rebalances
	go func() {
		for {
			if err := c.client.Consume(ctx, topics, consumer); err != nil {
				if c.errHandler != nil {
					c.errHandler(err)
				}
				time.Sleep(1 * time.Second)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

// Close closes the consumer group
func (c *consumerGroup) Close() error {
	return c.client.Close()
}
