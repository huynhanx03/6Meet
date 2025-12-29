package infrastructure

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/database/elasticsearch"
	"go.uber.org/zap"
)

// SetupElasticsearch initializes Elasticsearch connection
func SetupElasticsearch() {
	client, err := elasticsearch.New(global.Config.Elasticsearch)
	if err != nil {
		global.Logger.Fatal("Failed to setup Elasticsearch client", zap.Error(err))
	}

	global.ESClient = client
	global.Logger.Info("Connected to Elasticsearch successfully")
}
