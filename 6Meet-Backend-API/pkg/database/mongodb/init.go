package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/settings"
)

type Connection struct {
	Client *mongo.Client
	DB     *mongo.Database
	config *settings.MongoDBConfig
}

func (c *Connection) connect() error {
	if c.config.Timeout == 0 {
		c.config.Timeout = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.Timeout)*time.Second)
	defer cancel()

	uri := c.buildURI()

	clientOpts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(c.config.MaxPoolSize).
		SetMinPoolSize(c.config.MinPoolSize).
		SetMaxConnIdleTime(time.Duration(c.config.MaxConnIdleTime) * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	c.Client = client
	c.DB = client.Database(c.config.Database)

	return nil
}

func (c *Connection) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	return nil
}

func (c *Connection) buildURI() string {
	if c.config.Username != "" && c.config.Password != "" {
		return fmt.Sprintf(
			"mongodb://%s:%s@%s:%d",
			c.config.Username,
			c.config.Password,
			c.config.Host,
			c.config.Port,
		)
	}

	return fmt.Sprintf("mongodb://%s:%d", c.config.Host, c.config.Port)
}
