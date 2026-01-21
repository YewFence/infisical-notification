// Package repo 实现了 Repository 模式 (数据访问层)。
// 它封装了所有的数据库 CRUD 操作，将业务逻辑与底层 SQL/ORM 细节解耦。
package repo

import (
	"errors"
	"time"

	"backend/internal/models"

	"gorm.io/gorm"
)

// TodoRepository 结构体包含一个 GORM 数据库实例。
// 通过方法接收者 (receiver) 将数据库操作绑定到这个结构体上。
type TodoRepository struct {
	db *gorm.DB
}

// NewTodoRepository 创建并返回一个新的 Repository 实例。
// 这是一种常见的构造函数模式。
func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// List 返回所有待办事项，按 ID 倒序排列（最新的在前面）。
func (r *TodoRepository) List() ([]models.TodoItem, error) {
	var items []models.TodoItem
	// Find 方法会自动生成 SELECT * FROM todo_items 查询。
	// Order("id desc") 添加 ORDER BY id DESC 子句。
	// 结果被扫描到 items 切片中。
	if err := r.db.Order("id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Create 创建一个新的待办事项。
func (r *TodoRepository) Create(secretPath string, now time.Time) (models.TodoItem, error) {
	item := models.TodoItem{
		SecretPath:  secretPath,
		IsCompleted: false,
		CreatedAt:   now,
	}
	// Create 方法将结构体插入数据库。
	// 如果插入失败（例如违反唯一约束），会返回 error。
	if err := r.db.Create(&item).Error; err != nil {
		return models.TodoItem{}, err
	}
	return item, nil
}

// ToggleComplete 切换待办事项的完成状态。
// 如果当前为未完成，则标记为已完成并设置完成时间；
// 如果当前为已完成，则重置为未完成并清空完成时间。
func (r *TodoRepository) ToggleComplete(id uint, now time.Time) (models.TodoItem, error) {
	var item models.TodoItem
	// First 方法查找第一条匹配记录，如果没找到会返回 gorm.ErrRecordNotFound。
	if err := r.db.First(&item, id).Error; err != nil {
		return models.TodoItem{}, err
	}

	// 根据当前状态决定切换方向
	if item.IsCompleted {
		// 当前已完成 → 切换为未完成
		if err := r.db.Model(&item).Updates(map[string]interface{}{
			"is_completed": false,
			"completed_at": nil, // 清空完成时间
		}).Error; err != nil {
			return models.TodoItem{}, err
		}
		item.IsCompleted = false
		item.CompletedAt = nil
	} else {
		// 当前未完成 → 切换为已完成
		if err := r.db.Model(&item).Updates(map[string]interface{}{
			"is_completed": true,
			"completed_at": &now, // 设置完成时间
		}).Error; err != nil {
			return models.TodoItem{}, err
		}
		item.IsCompleted = true
		item.CompletedAt = &now
	}

	return item, nil
}

// Delete 根据 ID 删除待办事项。
func (r *TodoRepository) Delete(id uint) error {
	// Delete 方法生成 DELETE 语句。
	// 即使记录不存在，Delete 通常也不会报错。
	result := r.db.Delete(&models.TodoItem{}, id)
	if result.Error != nil {
		return result.Error
	}
	// 通过 RowsAffected 检查是否真的删除了数据。
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpsertFromWebhook 处理 Webhook 事件：如果记录存在则重置状态，如果不存在则创建。
// Upsert = Update + Insert
func (r *TodoRepository) UpsertFromWebhook(secretPath string, now time.Time) (models.TodoItem, error) {
	var item models.TodoItem
	// 尝试根据 secret_path 查找记录
	err := r.db.Where("secret_path = ?", secretPath).First(&item).Error

	if err == nil {
		// 1. 记录存在：重置为 "未完成" 状态。
		// 这意味着 Infisical 端发生了变更，需要重新处理这个 Todo。
		if err := r.db.Model(&item).Updates(map[string]interface{}{
			"is_completed": false,
			"completed_at": nil, // 将字段置为 NULL
		}).Error; err != nil {
			return models.TodoItem{}, err
		}
		item.IsCompleted = false
		item.CompletedAt = nil
		return item, nil
	}

	// 检查错误是否真的为 "记录未找到"
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.TodoItem{}, err // 发生了其他数据库错误
	}

	// 2. 记录不存在：创建新记录
	item = models.TodoItem{
		SecretPath:  secretPath,
		IsCompleted: false,
		CreatedAt:   now,
	}
	if err := r.db.Create(&item).Error; err != nil {
		return models.TodoItem{}, err
	}
	return item, nil
}
