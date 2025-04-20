package main

import (
	"fmt"
	"log"

	"github.com/akinoccc/web-tracing-admin/internal/api"
	"github.com/akinoccc/web-tracing-admin/internal/middleware"
	"github.com/akinoccc/web-tracing-admin/internal/model"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/akinoccc/web-tracing-admin/docs"
)

// @title Web Tracing Admin API
// @version 1.0
// @description Web Tracing SDK 错误监控后台 API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://github.com/akinoccc/web-tracing-admin/blob/main/LICENSE

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 初始化配置
	model.Setup()

	// 创建 Gin 实例
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.CORS())

	// 注册路由
	registerRoutes(r)

	// 启动服务器
	port := fmt.Sprintf(":%d", model.ServerSetting.HttpPort)
	log.Printf("Server started on http://localhost%s", port)
	r.Run(port)
}

// 注册路由
func registerRoutes(r *gin.Engine) {
	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 认证路由
	r.POST("/api/auth/login", api.Login)
	r.POST("/api/auth/register", api.Register)

	// 跟踪数据上报路由
	r.POST("/trackweb", api.TrackWeb)
	r.GET("/trackweb", api.TrackWeb)

	// 兼容示例项目的路由
	r.GET("/getAllTracingList", api.GetAllTracingList)
	r.POST("/cleanTracingList", api.CleanTracingList)
	r.GET("/getBaseInfo", api.GetBaseInfo)

	// 需要认证的路由
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.JWT())
	{
		// 项目路由
		apiGroup.POST("/projects", api.CreateProject)
		apiGroup.GET("/projects", api.GetProjects)
		apiGroup.GET("/projects/:id", api.GetProject)
		apiGroup.PUT("/projects/:id", api.UpdateProject)
		apiGroup.DELETE("/projects/:id", api.DeleteProject)

		// 事件路由
		apiGroup.GET("/events", api.GetEvents)
		apiGroup.GET("/events/:id", api.GetEventDetail)

		// 统计路由
		apiGroup.GET("/events/stats", api.GetErrorStats)
		apiGroup.GET("/events/stats/browser", api.GetErrorByBrowserStats)
		apiGroup.GET("/events/stats/os", api.GetErrorByOSStats)
		apiGroup.GET("/events/stats/device", api.GetErrorByDeviceStats)
		apiGroup.GET("/events/stats/error-type", api.GetErrorByTypeStats)
		apiGroup.GET("/events/stats/performance", api.GetPerformanceStats)
		apiGroup.GET("/events/stats/request-error", api.GetRequestErrorStats)
	}
}
