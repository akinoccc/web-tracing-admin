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

	// 数据上报接口 - 不需要认证
	r.POST("/api/trackweb", api.TrackWeb)

	// 认证路由
	r.POST("/api/auth/login", api.Login)
	r.POST("/api/auth/register", api.Register)

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

		// 错误监控路由
		apiGroup.GET("/errors", api.GetErrors)
		apiGroup.GET("/errors/:id", api.GetErrorDetail)
		apiGroup.GET("/errors/stats", api.GetErrorStats)

		// 性能监控路由
		apiGroup.GET("/performance", api.GetPerformance)
		apiGroup.GET("/performance/stats", api.GetPerformanceStats)
		apiGroup.GET("/performance/resources", api.GetResourcePerformance)

		// 用户行为路由
		apiGroup.GET("/behavior/pv", api.GetPageViews)
		apiGroup.GET("/behavior/clicks", api.GetClicks)
		apiGroup.GET("/behavior/stats", api.GetBehaviorStats)
	}
}
