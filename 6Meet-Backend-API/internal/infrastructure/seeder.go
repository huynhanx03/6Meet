package infrastructure

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/internal/crawl"
	"go.uber.org/zap"
)

// LoadData loads data
func LoadData() error {
	err := crawl.CrawlData("storages/seed_names.txt")

	if err != nil {
		global.Logger.Error("Failed to load data", zap.Error(err))
		return err
	}

	return nil
}
