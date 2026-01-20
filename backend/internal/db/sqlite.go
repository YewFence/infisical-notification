// Package db 负责数据库的底层连接和配置。
// 这里封装了 GORM 对 SQLite 的初始化逻辑。
package db

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open 初始化并返回一个 GORM 数据库连接实例。
// 它会自动创建数据库文件所在的目录，并配置连接池参数。
func Open(path string) (*gorm.DB, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("empty sqlite path")
	}

	// 1. 确保数据库目录存在
	// filepath.Dir 获取路径中的目录部分。
	// os.MkdirAll 类似于 mkdir -p，递归创建目录。
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	// 2. 打开数据库连接
	// gorm.Open 是 GORM 的入口方法。
	// sqlite.Open(path) 指定使用 SQLite 驱动。
	// &gorm.Config{} 可以传递高级配置，这里使用默认配置。
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 3. 配置连接池
	// 虽然 GORM 返回的是 *gorm.DB，但底层的 *sql.DB 对象提供了连接池配置能力。
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SQLite 是文件型数据库，默认只支持单写。
	// 为了避免 "database is locked" 错误，我们将最大打开连接数限制为 1。
	// 这意味着所有写操作会串行执行，读操作也会受到影响，但保证了数据安全性。
	// 对于高并发场景，通常会切换到 PostgreSQL 或 MySQL。
	sqlDB.SetMaxOpenConns(1)

	// 设置最大空闲连接数，保持 1 个连接常驻，避免频繁打开/关闭文件。
	sqlDB.SetMaxIdleConns(1)

	// 设置连接最大生命周期，防止连接长时间闲置导致的潜在问题（虽然 SQLite 不太涉及连接超时）。
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}