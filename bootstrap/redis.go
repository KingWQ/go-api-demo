package bootstrap

import (
	"fmt"
	"go-api-demo/pkg/config"
	"go-api-demo/pkg/redis"
)

// SetupRedis 初始化redis
func SetupRedis() {
	redis.ConnectRedis(
		fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		config.GetString("redis.username"),
		config.GetString("redis.password"),
		config.GetInt("redis.database"),
	)
}
