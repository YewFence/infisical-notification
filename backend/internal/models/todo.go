// TODO 数据模型定义，用于 GORM 映射。
package models

import "time"

// TodoItem represents a TODO entry stored in SQLite.
type TodoItem struct {
	ID          uint       `gorm:"primaryKey"`
	SecretPath  string     `gorm:"column:secret_path;uniqueIndex;not null"`
	IsCompleted bool       `gorm:"column:is_completed;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
}

// TableName keeps table naming stable for migrations.
func (TodoItem) TableName() string {
	return "todo_items"
}

