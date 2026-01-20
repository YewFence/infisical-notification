// 配置加载模块，集中处理环境变量与默认值。
package config

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultBindAddr = ":8080"
)

// Config holds backend configuration values.
type Config struct {
	WebhookSecret string
	DBPath        string
	BindAddr      string
}

// Load reads configuration from env with defaults.
func Load() (Config, error) {
	cfg := Config{
		WebhookSecret: strings.TrimSpace(os.Getenv("INFISICAL_WEBHOOK_SECRET")),
		DBPath:        strings.TrimSpace(os.Getenv("TODO_DB_PATH")),
		BindAddr:      strings.TrimSpace(os.Getenv("TODO_BIND_ADDR")),
	}

	if cfg.DBPath == "" {
		cfg.DBPath = defaultDBPath()
	}
	if cfg.BindAddr == "" {
		cfg.BindAddr = defaultBindAddr
	}

	return cfg, nil
}

func defaultDBPath() string {
	if _, err := os.Stat("backend"); err == nil {
		return filepath.Join("backend", "data", "todos.db")
	}
	return filepath.Join("data", "todos.db")
}

