// Package config 负责应用程序的配置管理。
// 它遵循 "配置即代码" 的原则，明确定义了应用所需的所有配置项。
package config

import (
	"os"
	"path/filepath"
	"strings"
)

// defaultBindAddr 定义默认的监听地址和端口。
// ":8080" 表示监听所有网卡的 8080 端口。
const (
	defaultBindAddr = ":8080"
)

// Config 结构体定义了所有可配置的参数。
// 使用结构体可以将分散的环境变量聚合在一起，方便在程序中传递。
type Config struct {
	// WebhookSecret 用于验证 Infisical 发来的 Webhook 请求的签名。
	// 这是一个敏感信息，必须通过环境变量注入。
	WebhookSecret string

	// DBPath 指定 SQLite 数据库文件的存储路径。
	DBPath string

	// BindAddr 指定 HTTP 服务监听的地址，例如 ":8080" 或 "127.0.0.1:3000"。
	BindAddr string
}

// Load 从环境变量加载配置，并应用默认值。
// 返回配置对象或错误。
func Load() (Config, error) {
	// 初始化配置对象，直接从 os.Getenv 读取环境变量。
	// strings.TrimSpace 用于去除可能存在的首尾空格，防止配置错误。
	cfg := Config{
		WebhookSecret: strings.TrimSpace(os.Getenv("INFISICAL_WEBHOOK_SECRET")),
		DBPath:        strings.TrimSpace(os.Getenv("TODO_DB_PATH")),
		BindAddr:      strings.TrimSpace(os.Getenv("TODO_BIND_ADDR")),
	}

	// 设置默认值逻辑
	// 如果环境变量未设置（空字符串），则使用预定义的默认值。

	if cfg.DBPath == "" {
		cfg.DBPath = defaultDBPath()
	}
	if cfg.BindAddr == "" {
		cfg.BindAddr = defaultBindAddr
	}

	return cfg, nil
}

// defaultDBPath 计算数据库的默认路径。
// 它会检查当前目录下是否存在 "backend" 文件夹，以适配不同的运行环境（项目根目录 vs backend 子目录）。
func defaultDBPath() string {
	// 如果存在 backend 目录，说明可能是在项目根目录运行，数据库放在 backend/data/todos.db
	if _, err := os.Stat("backend"); err == nil {
		return filepath.Join("backend", "data", "todos.db")
	}
	// 否则假设在 backend 目录下运行，数据库放在 data/todos.db
	return filepath.Join("data", "todos.db")
}