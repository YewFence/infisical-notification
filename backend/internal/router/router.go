// 路由装配与中间件初始化。
package router

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repo"

	"github.com/gin-gonic/gin"
)

// NewRouter wires the API routes.
func NewRouter(cfg config.Config, repo *repo.TodoRepository) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	todoHandler := handlers.NewTodoHandler(repo)
	webhookHandler := handlers.NewWebhookHandler(repo, cfg.WebhookSecret)

	api := engine.Group("/api/todos")
	api.POST("/webhook", webhookHandler.Handle)
	api.GET("", todoHandler.List)
	api.POST("", todoHandler.Create)
	api.PATCH("/:id", todoHandler.Update)
	api.DELETE("/:id", todoHandler.Delete)
	api.POST("/:id/complete", todoHandler.Complete)

	return engine
}

