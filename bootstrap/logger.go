package bootstrap

import (
	"go-api-demo/pkg/config"
	"go-api-demo/pkg/logger"
)

// SetupLogger 初始化Logger
func SetupLogger() {
	logger.InitLogger(
		config.GetString("log.filename"),
		config.GetInt("log.max_size"),
		config.GetInt("log.max_backup"),
		config.GetInt("log.max_age"),
		config.GetBool("log.compress"),
		config.GetString("log.type"),
		config.GetString("log.level"),
	)
}
