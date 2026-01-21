package api

import (
	"github.com/gin-gonic/gin"
)

// 路由注册
func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// 用于健康检查的 ping 接口
		api.GET("/ping", PingHandler) // 我们需要定义一个 PingHandler
		// 用于分析布尔函数性质的接口
		api.POST("/analyze", AnalyzeFunctionHandler)
	}
}
