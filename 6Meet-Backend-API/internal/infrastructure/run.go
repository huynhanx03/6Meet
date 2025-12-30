package infrastructure

import (
	"context"

	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/di"

	"go.uber.org/zap"
)

func Run() error {
	LoadConfig()

	SetupLogger()
	SetupMongoDB()
	SetupRedis()
	SetupKafka()
	SetupElasticsearch()
	di.SetupDependencies()
	http := NewHTTPServer()

	ctx := context.Background()

	if err := di.GlobalContainer.UserCDCConsumer.Start(ctx); err != nil {
		global.Logger.Error("CDC Consumer failed", zap.Error(err))
	}

	if err := di.GlobalContainer.CrawlService.StartCrawling(ctx, "storages/seed_names.txt"); err != nil {
		global.Logger.Error("Crawler failed", zap.Error(err))
	}

	return http.Run()
}
