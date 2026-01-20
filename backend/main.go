// TODO 后端服务入口，负责初始化配置、数据库和路由。
package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/repo"
	"backend/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.AutoMigrate(&models.TodoItem{}); err != nil {
		log.Fatal(err)
	}

	todoRepo := repo.NewTodoRepository(database)
	engine := router.NewRouter(cfg, todoRepo)
	if err := engine.Run(cfg.BindAddr); err != nil {
		log.Fatal(err)
	}
}

