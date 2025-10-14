package mongodb

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/settings"
)

func NewConnection(config *settings.MongoDBConfig) (*Connection, error) {
	conn := &Connection{
		config: config,
	}

	if err := conn.connect(); err != nil {
		return nil, err
	}

	return conn, nil
}
