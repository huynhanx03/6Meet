package di

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/cache"
	db "github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/db"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/producer"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driven/search"
	userconsumer "github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driver/consumer/user"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/adapters/driver/http"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/core/service"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/ports"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/mq/kafka"
	"go.uber.org/zap"
)

// InitUserDependencies initializes user related dependencies
func InitUserDependencies() (ports.UserRepository, ports.UserCacheRepository, ports.UserSearchRepository, ports.UserService, http.UserHandler, ports.CrawlService, ports.UserProducer, ports.UserConsumer, ports.UserConsumer) {
	// Initialize repositories
	userRepo := db.NewUser(global.MongoDB.DB)

	// Initialize producer
	userProducer := producer.NewUser(global.KafkaProducer)

	// Initialize consumer
	kafkaCfg := &kafka.Config{
		Brokers:  global.Config.Kafka.Brokers,
		ClientID: "wiki-crawler",
		ProducerInfo: kafka.ProducerConfig{
			FlushFrequency:  global.Config.Kafka.FlushFrequency,
			FlushBytes:      global.Config.Kafka.FlushBytes,
			MaxMessageBytes: global.Config.Kafka.MaxMessageBytes,
			MaxRetries:      global.Config.Kafka.MaxRetries,
			RetryBackoff:    global.Config.Kafka.RetryBackoff,
			ReturnSuccesses: true,
		},
		ConsumerInfo: kafka.ConsumerConfig{
			SessionTimeout:    global.Config.Kafka.Timeout * 1000,
			MaxProcessingTime: global.Config.Kafka.MaxProcessingTime,
		},
	}

	crawlConsumer, err := userconsumer.NewCrawlerConsumer(kafkaCfg, userRepo)
	if err != nil {
		global.Logger.Fatal("failed to create user consumer", zap.Error(err))
	}

	// Initialize cache
	userCache := cache.NewUser(global.Redis)

	// Initialize search
	userSearch := search.NewUser(global.ESClient)

	// Initialize services
	userService := service.NewUser(userRepo, userCache, userSearch)

	// Initialize handlers
	userHandler := http.NewUser(userService)

	crawlService, err := service.NewCrawlService(userProducer, crawlConsumer)
	if err != nil {
		global.Logger.Fatal("failed to create crawl service", zap.Error(err))
	}

	// Consumer now uses userService directly
	cdcConsumer, err := userconsumer.NewCDCConsumer(kafkaCfg, userService)
	if err != nil {
		global.Logger.Fatal("failed to create user cdc consumer", zap.Error(err))
	}

	return userRepo, userCache, userSearch, userService, userHandler, crawlService, userProducer, crawlConsumer, cdcConsumer
}
