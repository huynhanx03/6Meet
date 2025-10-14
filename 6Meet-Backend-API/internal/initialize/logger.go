package initialize

import (
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/global"
	"github.com/huynhanx03/6Meet/6Meet-Backend-API/pkg/logger"
)

// InitLogger initializes the logger
func InitLogger() {
	config := logger.LoggerConfig{
		Level:      global.Config.Logger.LogLevel,
		Filename:   global.Config.Logger.FileLogName,
		MaxSize:    global.Config.Logger.MaxSize,
		MaxBackups: global.Config.Logger.MaxBackups,
		MaxAge:     global.Config.Logger.MaxAge,
		Compress:   global.Config.Logger.Compress,
	}

	global.Logger = logger.NewLogger(config)
}
