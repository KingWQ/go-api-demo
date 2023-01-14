// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
	"github.com/gin-gonic/gin"
	"go-api-demo/app/http/middlewares"
	"go-api-demo/routes"
	"net/http"
	"strings"
)

// SetupRoute 路由初始化
func SetupRoute(router *gin.Engine) {
	//注册全局中间件
	registerGlobalMiddleWare(router)

	//注册API路由
	routes.RegisterAPIRoutes(router)

	//配置 404 路由
	setup404Handler(router)
}

func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(middlewares.Logger(), middlewares.Recovery())
}

func setup404Handler(router *gin.Engine) {
	//处理 404 请求
	router.NoRoute(func(c *gin.Context) {
		//获取表头信息的Accept信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			c.String(http.StatusNotFound, "页面返回 404")
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认 url和请求方法是否正确。",
			})
		}
	})
}
