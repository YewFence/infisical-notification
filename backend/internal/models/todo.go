// Package models 定义了应用程序的数据模型。
// 这些结构体对应数据库中的表结构 (Schema)。
package models

import "time"

// TodoItem 代表待办事项表中的一行记录。
// GORM 使用结构体标签 (Tag) 来定义列的元数据（如主键、索引、约束等）。
type TodoItem struct {
	// ID 是主键，GORM 默认使用 uint 类型的 ID 字段作为主键。
	// `gorm:"primaryKey"` 显式声明这是主键。
	ID uint `gorm:"primaryKey"`

	// SecretPath 存储密钥路径，例如 "/dev/db/password"。
	// `gorm:"column:secret_path"` 指定数据库中的列名为 secret_path。
	// `uniqueIndex` 创建唯一索引，确保同一个路径只能有一条记录。
	// `not null` 约束该字段不能为空。
	SecretPath string `gorm:"column:secret_path;uniqueIndex;not null"`

	// IsCompleted 标记该待办事项是否已完成。
	// 使用 bool 类型，SQLite 中会存储为 0 或 1。
	IsCompleted bool `gorm:"column:is_completed;not null"`

	// CreatedAt 记录创建时间。
	// GORM 约定：如果字段名为 CreatedAt，它会在创建记录时自动填充当前时间。
	CreatedAt time.Time `gorm:"column:created_at;not null"`

	// CompletedAt 记录完成时间。
	// 使用指针类型 *time.Time 是为了支持 NULL 值。
	// 如果该字段是 nil，数据库中存储为 NULL，表示尚未完成。
	CompletedAt *time.Time `gorm:"column:completed_at"`
}

// TableName 实现 GORM 的 Tabler 接口，用于自定义表名。
// 如果不实现此方法，GORM 默认会将结构体名称转换为复数蛇形命名（如 todo_items）。
// 这里显式返回 "todo_items" 只是为了明确性。
func (TodoItem) TableName() string {
	return "todo_items"
}