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

// todoInput 定义了创建接口的请求体结构。
type todoInput struct {
	// 反射标签指定了 JSON 字段名。
	SecretPath string `json:"secretPath"`
}

// List 获取所有待办事项列表。
//
//	@Summary		获取待办事项列表
//	@Description	获取所有待办事项的列表
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"成功返回待办事项列表"
//	@Failure		500	{object}	map[string]string		"服务器内部错误"
//	@Router			/ [get]
func (h *TodoHandler) List(c *gin.Context) {
	items, err := h.repo.List()
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "list todos failed")
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
//
//	@Summary		创建待办事项
//	@Description	创建一个新的待办事项
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			todo	body		todoInput				true	"待办事项信息"
//	@Success		200		{object}	map[string]interface{}	"成功返回创建的待办事项"
//	@Failure		400		{object}	map[string]string		"请求参数错误"
//	@Failure		409		{object}	map[string]string		"密钥路径已存在"
//	@Failure		500		{object}	map[string]string		"服务器内部错误"
//	@Router			/ [post]
func (h *TodoHandler) Create(c *gin.Context) {
	var input todoInput
	// ShouldBindJSON 解析请求体中的 JSON 并绑定到 input 结构体。
	// 如果 JSON 格式错误或字段类型不匹配，返回 error。
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	secretPath := strings.TrimSpace(input.SecretPath)
	if secretPath == "" {
		RespondError(c, http.StatusBadRequest, "secretPath is required")
		return
	}

	item, err := h.repo.Create(secretPath, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			RespondError(c, http.StatusConflict, "secretPath already exists")
			return
		}
		RespondError(c, http.StatusInternalServerError, "create todo failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

// Get 获取单个待办事项。
//
//	@Summary		获取单个待办事项
//	@Description	根据 ID 获取单个待办事项的详细信息
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"待办事项 ID"
//	@Success		200		{object}	map[string]interface{}	"成功返回待办事项"
//	@Failure		400		{object}	map[string]string		"请求参数错误"
//	@Failure		404		{object}	map[string]string		"待办事项不存在"
//	@Failure		500		{object}	map[string]string		"服务器内部错误"
//	@Router			/{id} [get]
func (h *TodoHandler) Get(c *gin.Context) {
	// 未找到 ID 则返回 400 错误
	id, ok := parseID(c)
	if !ok {
		return
	}

	item, err := h.repo.GetByID(id)
	if err != nil {
		// 处理未找到的情况
		if errors.Is(err, gorm.ErrRecordNotFound) {
			RespondError(c, http.StatusNotFound, "todo not found")
			return
		}
		// 其他错误
		RespondError(c, http.StatusInternalServerError, "get todo failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

// ToggleComplete 切换待办事项的完成状态。
//
//	@Summary		切换待办事项完成状态
//	@Description	切换指定 ID 的待办事项的完成状态（已完成↔未完成）
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"待办事项 ID"
//	@Success		200		{object}	map[string]interface{}	"成功返回切换后的待办事项"
//	@Failure		400		{object}	map[string]string		"请求参数错误"
//	@Failure		404		{object}	map[string]string		"待办事项不存在"
//	@Failure		500		{object}	map[string]string		"服务器内部错误"
//	@Router			/{id} [patch]
func (h *TodoHandler) ToggleComplete(c *gin.Context) {
	// 从 URL 参数获取 ID
	// 未找到 ID 则返回 400 错误
	id, ok := parseID(c)
	if !ok {
		return
	}

	item, err := h.repo.ToggleComplete(id, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			RespondError(c, http.StatusNotFound, "todo not found")
			return
		}
		RespondError(c, http.StatusInternalServerError, "toggle complete failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

// Delete 删除待办事项。
//
//	@Summary		删除待办事项
//	@Description	删除指定 ID 的待办事项
//	@Tags			todos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"待办事项 ID"
//	@Success		200	{object}	map[string]string    	"成功删除"
//	@Failure		400	{object}	map[string]string		"请求参数错误"
//	@Failure		404	{object}	map[string]string		"待办事项不存在"
//	@Failure		500	{object}	map[string]string		"服务器内部错误"
//	@Router			/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			RespondError(c, http.StatusNotFound, "todo not found")
			return
		}
		RespondError(c, http.StatusInternalServerError, "delete todo failed")
		return
	}

	respondOK(c, "ok")
}

// parseID 辅助函数：从 URL 路径参数中解析 uint 类型的 ID。
// 示例：/api/todos/123 -> 123
func parseID(c *gin.Context) (uint, bool) {
	idText := strings.TrimSpace(c.Param("id"))
	idValue, err := strconv.ParseUint(idText, 10, 64)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(idValue), true
}
