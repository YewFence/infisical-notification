// Package config 负责应用程序的配置管理。
// 它遵循 "配置即代码" 的原则，明确定义了应用所需的所有配置项。
package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// defaultBindAddr 定义默认的监听地址和端口。
// ":8080" 表示监听所有网卡的 8080 端口。
// defaultMaxBodySize 定义默认的请求体大小限制（10MB）。
const (
	defaultBindAddr     = ":8080"
	defaultMaxBodySize  = 10 << 20 // 10MB
)

// Config 结构体定义了所有可配置的参数。
// 使用结构体可以将分散的环境变量聚合在一起，方便在程序中传递。
type Config struct {
	// Environment 指定运行环境。
	// 支持的值：development/dev（开发环境）、production/prod（生产环境）。
	// 默认为 development。
	Environment string

	// WebhookSecret 用于验证 Infisical 发来的 Webhook 请求的签名。
	// 这是一个敏感信息，必须通过环境变量注入。
	WebhookSecret string

	// DBPath 指定 SQLite 数据库文件的存储路径。
	DBPath string

	// BindAddr 指定 HTTP 服务监听的地址，例如 ":8080" 或 "127.0.0.1:3000"。
	BindAddr string

	// MaxBodySize 指定请求体的最大大小（字节）。
	MaxBodySize int64

	// CORSAllowedOrigins 指定允许的跨域来源列表。
	// 开发环境：为空时允许 localhost 和 127.0.0.1 的所有端口。
	// 生产环境：应设置具体的域名，多个域名用逗号分隔。
	CORSAllowedOrigins []string
}

// IsDevelopment 判断是否为开发模式。
// 支持 development/dev 或空值（默认为开发环境）。
func (c *Config) IsDevelopment() bool {
	env := strings.ToLower(c.Environment)
	return env == "" || env == "development" || env == "dev"
}

// IsProduction 判断是否为生产模式。
// 支持 production/prod。
func (c *Config) IsProduction() bool {
	env := strings.ToLower(c.Environment)
	return env == "production" || env == "prod"
}

// Load 从环境变量加载配置，并应用默认值。
// 返回配置对象或错误。
func Load() (Config, error) {
	// 初始化配置对象，直接从 os.Getenv 读取环境变量。
	// strings.TrimSpace 用于去除可能存在的首尾空格，防止配置错误。
	cfg := Config{
		Environment:   strings.TrimSpace(os.Getenv("APP_ENV")),
		WebhookSecret: strings.TrimSpace(os.Getenv("INFISICAL_WEBHOOK_SECRET")),
		DBPath:        strings.TrimSpace(os.Getenv("TODO_DB_PATH")),
		BindAddr:      strings.TrimSpace(os.Getenv("TODO_BIND_ADDR")),
	}

	// 设置默认值逻辑
	// 如果环境变量未设置（空字符串），则使用预定义的默认值。

	// 环境变量未设置时输出警告
	if cfg.Environment == "" {
		println("警告: APP_ENV 环境变量未设置，默认使用开发环境 (development)")
	}

	if cfg.DBPath == "" {
		cfg.DBPath = defaultDBPath()
	}
	if cfg.BindAddr == "" {
		cfg.BindAddr = defaultBindAddr
	}

	// 加载请求体大小限制配置
	maxBodySizeStr := strings.TrimSpace(os.Getenv("TODO_MAX_BODY_SIZE"))
	if maxBodySizeStr != "" {
		if size, err := strconv.ParseInt(maxBodySizeStr, 10, 64); err == nil && size > 0 {
			cfg.MaxBodySize = size
		} else {
			cfg.MaxBodySize = defaultMaxBodySize
		}
	} else {
		cfg.MaxBodySize = defaultMaxBodySize
	}

	// 加载 CORS 配置
	corsOrigins := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if corsOrigins != "" {
		origins := strings.Split(corsOrigins, ",")
		for _, origin := range origins {
			if trimmed := strings.TrimSpace(origin); trimmed != "" {
				cfg.CORSAllowedOrigins = append(cfg.CORSAllowedOrigins, trimmed)
			}
		}
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