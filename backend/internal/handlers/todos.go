// Package handlers 实现了具体的业务逻辑处理器 (Controller)。
// 这里处理与 Todo 相关的 HTTP 请求。
package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/internal/repo"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TodoHandler 结构体持有 Repository 的引用。
// 这样可以在处理请求时调用数据库操作。
type TodoHandler struct {
	repo *repo.TodoRepository
}

// NewTodoHandler 创建一个新的 TodoHandler。
func NewTodoHandler(repo *repo.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

// todoInput 定义了创建/更新接口的请求体结构。
type todoInput struct {
	// 反射标签指定了 JSON 字段名。
	SecretPath string `json:"secretPath"`
}

// List 获取所有待办事项列表。
// GET /api/todos
func (h *TodoHandler) List(c *gin.Context) {
	items, err := h.repo.List()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list todos failed")
		return
	}

	// 将数据库模型切片转换为响应模型切片
	response := make([]TodoResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toTodoResponse(item))
	}
	respondOK(c, response)
}

// Create 创建新的待办事项。
// POST /api/todos
func (h *TodoHandler) Create(c *gin.Context) {
	var input todoInput
	// ShouldBindJSON 解析请求体中的 JSON 并绑定到 input 结构体。
	// 如果 JSON 格式错误或字段类型不匹配，返回 error。
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	secretPath := strings.TrimSpace(input.SecretPath)
	if secretPath == "" {
		respondError(c, http.StatusBadRequest, "secretPath is required")
		return
	}

	item, err := h.repo.Create(secretPath, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			respondError(c, http.StatusConflict, "secretPath already exists")
			return
		}
		respondError(c, http.StatusInternalServerError, "create todo failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

// Update 更新待办事项的内容 (SecretPath)。
// PATCH /api/todos/:id
func (h *TodoHandler) Update(c *gin.Context) {
	// 从 URL 参数获取 ID
	id, ok := parseID(c)
	if !ok {
		return
	}

	var input todoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	secretPath := strings.TrimSpace(input.SecretPath)
	if secretPath == "" {
		respondError(c, http.StatusBadRequest, "secretPath is required")
		return
	}

	item, err := h.repo.UpdateSecretPath(id, secretPath)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "todo not found")
			return
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			respondError(c, http.StatusConflict, "secretPath already exists")
			return
		}
		respondError(c, http.StatusInternalServerError, "update todo failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

// Delete 删除待办事项。
// DELETE /api/todos/:id
func (h *TodoHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "todo not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "delete todo failed")
		return
	}

	respondOK(c, true)
}

// Complete 将待办事项标记为完成。
// POST /api/todos/:id/complete
func (h *TodoHandler) Complete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	item, err := h.repo.Complete(id, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(c, http.StatusNotFound, "todo not found")
			return
		}
		respondError(c, http.StatusInternalServerError, "complete todo failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

// parseID 辅助函数：从 URL 路径参数中解析 uint 类型的 ID。
// 示例：/api/todos/123 -> 123
func parseID(c *gin.Context) (uint, bool) {
	idText := strings.TrimSpace(c.Param("id"))
	idValue, err := strconv.ParseUint(idText, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(idValue), true
}
