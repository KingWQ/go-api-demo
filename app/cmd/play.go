package cmd

import (
	"github.com/spf13/cobra"
	"go-api-demo/pkg/console"
	"go-api-demo/pkg/redis"
	"time"
)

var CmdPlay = &cobra.Command{
	Use:   "play",
	Short: "Likes the Go Playground, but running at our application context",
	Run:   runPlay,
}

// 调试完成后请记得清除测试代码
func runPlay(cmd *cobra.Command, args []string) {
	redis.Redis.Set("hello", "hi from redis", 10*time.Second)
	console.Success(redis.Redis.Get("hello"))
}
