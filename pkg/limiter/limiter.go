// Package limiter 处理限流逻辑
package limiter

import (
	"github.com/gin-gonic/gin"
	limiterlib "github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"go-api-demo/pkg/config"
	"go-api-demo/pkg/redis"
	"golang.org/x/net/context"

	"go-api-demo/pkg/logger"
	"strings"
)

// GetKeyIP 获取Limiter 的key, ip
func GetKeyIP(c *gin.Context) string {
	return c.ClientIP()
}

// GetKeyRouteWithIP Limiter 的key, 路由+IP, 针对单个路由做限流
func GetKeyRouteWithIP(c *gin.Context) string {
	return routeToKeyString(c.FullPath()) + c.ClientIP()
}

// CheckRate 检测请求是否超额
func CheckRate(c *gin.Context, key string, formatted string) (limiterlib.Context, error) {
	// 1. 实例化依赖的 limiter 包的 limiter.Rate 对象
	var contextLimit limiterlib.Context
	rate, err := limiterlib.NewRateFromFormatted(formatted)
	if err != nil {
		logger.LogIf(err)
		return contextLimit, err
	}

	// 2. 初始化存储，使用我们程序里共用的 redis.Redis 对象
	store, err := sredis.NewStoreWithOptions(redis.Redis.Client, limiterlib.StoreOptions{
		// 为 limiter 设置前缀，保持 redis 里数据的整洁
		Prefix: config.GetString("app.name") + ":limiter",
	})
	if err != nil {
		logger.LogIf(err)
		return contextLimit, err
	}

	// 3. 使用上面的初始化的 limiter.Rate 对象和存储对象
	limiterObj := limiterlib.New(store, rate)

	// 4. 获取限流的结果
	ctx := context.WithValue(context.Background(), "GinContext", c)
	if c.GetBool("limiter-once") {
		// Peek() 取结果，不增加访问次数
		return limiterObj.Peek(ctx, key)
	} else {

		// 确保多个路由组里调用 LimitIP 进行限流时，只增加一次访问次数。
		c.Set("limiter-once", true)

		// Get() 取结果且增加访问次数
		return limiterObj.Get(ctx, key)
	}
}

// RouteToKeyString 辅助方法，将URL中的 / 格式 为 -
func routeToKeyString(routeName string) string {
	routeName = strings.ReplaceAll(routeName, "/", "-")
	routeName = strings.ReplaceAll(routeName, ":", "_")

	return routeName
}
