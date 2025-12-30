package producer

import (
	"context"
	"errors"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/constant"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/event"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/mq/kafka"
)

var (
	// ErrInvalidEventType is returned when the event type is invalid
	ErrInvalidEventType = errors.New("invalid event type")
)

type userProducer struct {
	producer kafka.Producer
}

// NewUser creates a new user producer
func NewUser(producer kafka.Producer) ports.UserProducer {
	return &userProducer{
		producer: producer,
	}
}

// Publish publishes a user related event to Kafka
func (p *userProducer) Publish(ctx context.Context, evt any) error {
	switch e := evt.(type) {
	case *event.UserCrawled:
		return kafka.PublishJSON(ctx, p.producer, constant.TopicCrawlWikiPages, func(data *event.UserCrawled) string {
			return data.Name
		}, e)
	default:
		return ErrInvalidEventType
	}
}
