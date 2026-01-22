// Package router 负责路由配置和 HTTP 服务初始化。
// 类似于 Flask 的 app.route 或 Django 的 urls.py。
package router

import (
	"strings"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repo"

	_ "backend/docs" // 导入生成的 Swagger 文档

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter 构造并配置 Gin 引擎。
// 这里进行了依赖注入：将 repo 注入到 handlers，再将 handlers 注册到路由。
func NewRouter(cfg config.Config, repo *repo.TodoRepository) *gin.Engine {
	// gin.New() 创建一个空白的 Gin 实例，不包含任何默认中间件。
	// 相比 gin.Default()，这给了我们更多的定制空间。
	engine := gin.New()

	// 注册全局中间件：
	// gin.Logger(): 将请求日志输出到控制台。
	// gin.Recovery(): 捕获任何 panic，防止程序崩溃，并返回 500 错误。
	engine.Use(gin.Logger(), gin.Recovery())

	// 添加请求体大小限制中间件
	engine.Use(middleware.BodySizeLimit(cfg.MaxBodySize))

	// 配置 CORS 中间件，允许前端跨域访问
	engine.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// 允许 localhost 和 127.0.0.1 的所有端口
			return strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:")
		},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 初始化业务处理器 (Handlers)
	todoHandler := handlers.NewTodoHandler(repo)
	webhookHandler := handlers.NewWebhookHandler(repo, cfg.WebhookSecret)

	// 创建路由组 (Route Group)
	// 所有以 /api/todos 开头的请求都会进入这个分组。
	api := engine.Group("/api/todos")
	{
		// 注册具体的路由规则：

		// Webhook 接口，用于接收外部系统 (Infisical) 的通知
		api.POST("/webhook", webhookHandler.Handle)

		// 标准 RESTful 接口
		api.GET("", todoHandler.List)                 // 获取列表
		api.POST("", todoHandler.Create)              // 创建
		api.GET("/:id", todoHandler.Get)              // 获取单个待办事项
		api.PATCH("/:id", todoHandler.ToggleComplete) // 切换完成状态
		api.DELETE("/:id", todoHandler.Delete)        // 删除
	}

	// 注册 Swagger UI 路由
	// 访问 http://localhost:8080/swagger/index.html 查看 API 文档
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 提供 swagger.yaml 文件的直接访问
	// 访问 http://localhost:8080/swagger.yaml 获取 OpenAPI 规范文件
	// 可以将这个 URL 导入到 Postman、Apifox 等工具中
	engine.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	engine.StaticFile("/swagger.json", "./docs/swagger.json")

	return engine
}
