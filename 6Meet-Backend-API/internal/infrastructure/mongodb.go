package infrastructure

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/mongodb"
)

// SetupMongoDB initializes the MongoDB connection
func SetupMongoDB() {
	config := global.Config.MongoDB

	conn, err := mongodb.New(&config)

	if err != nil {
		global.Logger.Sugar().Fatalf("Failed to connect to MongoDB: %v", err)
	}

	global.MongoDB = conn
	global.Logger.Sugar().Info("Connected to MongoDB successfully")
}
