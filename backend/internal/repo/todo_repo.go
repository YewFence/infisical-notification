// TODO 数据访问层，封装数据库读写。
package repo

import (
	"errors"
	"time"

	"backend/internal/models"

	"gorm.io/gorm"
)

// TodoRepository wraps database operations for todo items.
type TodoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new repository.
func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// List returns all todo items ordered by id desc.
func (r *TodoRepository) List() ([]models.TodoItem, error) {
	var items []models.TodoItem
	if err := r.db.Order("id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Create inserts a new todo item.
func (r *TodoRepository) Create(secretPath string, now time.Time) (models.TodoItem, error) {
	item := models.TodoItem{
		SecretPath:  secretPath,
		IsCompleted: false,
		CreatedAt:   now,
	}
	if err := r.db.Create(&item).Error; err != nil {
		return models.TodoItem{}, err
	}
	return item, nil
}

// UpdateSecretPath updates secret path by id.
func (r *TodoRepository) UpdateSecretPath(id uint, secretPath string) (models.TodoItem, error) {
	var item models.TodoItem
	if err := r.db.First(&item, id).Error; err != nil {
		return models.TodoItem{}, err
	}

	if err := r.db.Model(&item).Updates(map[string]interface{}{
		"secret_path": secretPath,
	}).Error; err != nil {
		return models.TodoItem{}, err
	}
	item.SecretPath = secretPath
	return item, nil
}

// Delete removes a todo item by id.
func (r *TodoRepository) Delete(id uint) error {
	result := r.db.Delete(&models.TodoItem{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Complete marks a todo item as completed.
func (r *TodoRepository) Complete(id uint, now time.Time) (models.TodoItem, error) {
	var item models.TodoItem
	if err := r.db.First(&item, id).Error; err != nil {
		return models.TodoItem{}, err
	}

	if err := r.db.Model(&item).Updates(map[string]interface{}{
		"is_completed": true,
		"completed_at": &now,
	}).Error; err != nil {
		return models.TodoItem{}, err
	}

	item.IsCompleted = true
	item.CompletedAt = &now
	return item, nil
}

// UpsertFromWebhook inserts or resets a todo item by secret path.
func (r *TodoRepository) UpsertFromWebhook(secretPath string, now time.Time) (models.TodoItem, error) {
	var item models.TodoItem
	err := r.db.Where("secret_path = ?", secretPath).First(&item).Error
	if err == nil {
		if err := r.db.Model(&item).Updates(map[string]interface{}{
			"is_completed": false,
			"completed_at": nil,
		}).Error; err != nil {
			return models.TodoItem{}, err
		}
		item.IsCompleted = false
		item.CompletedAt = nil
		return item, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.TodoItem{}, err
	}

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

