package infrastructure

import (
	"log"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/mq/kafka"
)

const (
	clientID = "6meet-backend-producer"
)

// SetupKafka initializes the Kafka producer
func SetupKafka() {
	if len(global.Config.Kafka.Brokers) == 0 {
		log.Println("Kafka configuration is missing, skipping setup")
		return
	}

	cfg := &kafka.Config{
		Brokers:  global.Config.Kafka.Brokers,
		ClientID: clientID,
		ProducerInfo: kafka.ProducerConfig{
			FlushFrequency:  global.Config.Kafka.FlushFrequency,
			FlushBytes:      global.Config.Kafka.FlushBytes,
			MaxMessageBytes: global.Config.Kafka.MaxMessageBytes,
			MaxRetries:      global.Config.Kafka.MaxRetries,
			RetryBackoff:    global.Config.Kafka.RetryBackoff,
			ReturnSuccesses: true,
		},
	}

	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	global.KafkaProducer = producer
	global.Logger.Info("Kafka Producer initialized successfully")
}
