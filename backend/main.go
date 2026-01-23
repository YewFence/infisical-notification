// Package main 是整个后端服务的入口。
// 在 Go 语言中,main 包是可执行程序的起点。
//
//	@title						Infisical Notification API
//	@version					1.0
//	@description				这是一个接收 Infisical Webhook 通知并管理 Todo 任务的后端服务
//	@termsOfService				http://swagger.io/terms/
//
//	@contact.name				API Support
//	@contact.email				support@example.com
//
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//
//	@host						localhost:8080
//	@BasePath					/api/todos
//	@schemes					http https

package main

import (
	"log"
	"log/slog"
	"strings"

	"github.com/joho/godotenv"

	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/repo"
	"backend/internal/router"
)

// main 函数是程序的执行入口,类似于 Python 的 if __name__ == "__main__": 下的代码。
// 这里负责组装各个组件(配置、数据库、路由)并启动 Web 服务。
func main() {
	// 0. 加载 .env 文件(如果存在)
	// 在开发环境中,可以通过 .env 文件设置环境变量。
	// 生产环境通常直接使用系统环境变量,所以这里忽略文件不存在的错误。
	_ = godotenv.Load()

	// 1. 加载配置
	// 从环境变量中读取配置信息,如果未设置则使用默认值。
	// 这符合 "12-Factor App" 的推荐实践。
	cfg, err := config.Load()
	if err != nil {
		// log.Fatal 会打印错误日志并以非零状态码退出程序 (os.Exit(1))。
		log.Fatal(err)
	}

	// 打印启动配置信息，便于排查问题
	corsDisplay := "(未配置，开发模式允许 localhost)"
	if len(cfg.CORSAllowedOrigins) > 0 {
		corsDisplay = strings.Join(cfg.CORSAllowedOrigins, ", ")
	}
	slog.Info("服务配置已加载",
		"environment", cfg.Environment,
		"bind_addr", cfg.BindAddr,
		"cors_origins", corsDisplay,
	)

	// 2. 初始化数据库连接
	// 使用 GORM 连接 SQLite 数据库。
	// 这里会处理数据库文件的创建和连接池的设置。
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 自动迁移 (Auto Migration)
	// GORM 的一个强大功能,它会根据 Go 的结构体定义自动创建或更新数据库表结构。
	// 类似于 Django 的 makemigrations/migrate 或 Flask-Migrate,但它是运行时自动完成的。
	// 这里确保 todo_items 表存在且字段正确。
	if err := database.AutoMigrate(&models.TodoItem{}); err != nil {
		log.Fatal(err)
	}

	// 4. 初始化 Repository (数据访问层)
	// 将数据库连接注入到 Repository 中。所有数据库操作都通过 todoRepo 进行。
	todoRepo := repo.NewTodoRepository(database)

	// 5. 初始化 Router (路由层)
	// 将配置和 Repository 注入到 Router 中。
	// Router 负责设置 HTTP 路由规则,并将请求分发给对应的 Handler。
	engine := router.NewRouter(cfg, todoRepo)

	// 6. 启动 Web 服务
	// Run() 方法会监听指定的端口(例如 :8080)并开始处理请求。
	// 这不仅会阻塞当前 goroutine,还会监听中断信号以优雅关闭(虽然 Gin 默认 Run 实现比较简单,生产环境可能需要更复杂的优雅关闭逻辑)。
	if err := engine.Run(cfg.BindAddr); err != nil {
		log.Fatal(err)
	}
}
