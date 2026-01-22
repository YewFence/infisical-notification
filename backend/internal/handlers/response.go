// Package handlers 定义了 HTTP 请求的响应格式和辅助函数。
// 这里旨在统一 API 的返回结构，使其符合 RESTful 风格或特定的前后端约定。
package handlers

import (
	"log"
	"net/http"
	"time"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

// TodoResponse 定义了返回给前端的 JSON 结构。
// 使用 `json:"..."` 标签控制序列化时的字段名。
// 前后端分离开发中，通常返回驼峰命名 (camelCase) 的 JSON 字段。
type TodoResponse struct {
	ID          uint       `json:"id"`
	SecretPath  string     `json:"secretPath"`
	IsCompleted bool       `json:"isCompleted"`
	CreatedAt   string     `json:"createdAt"`   // 格式化后的时间字符串
	CompletedAt *string    `json:"completedAt"` // 指针类型，允许为 null
}

const timeLayout = time.RFC3339

// toTodoResponse 将数据库模型转换为 API 响应模型。
// 这种 DTO (Data Transfer Object) 模式可以隔离数据库结构和 API 契约。
func toTodoResponse(item models.TodoItem) TodoResponse {
	response := TodoResponse{
		ID:          item.ID,
		SecretPath:  item.SecretPath,
		IsCompleted: item.IsCompleted,
		CreatedAt:   item.CreatedAt.Format(timeLayout),
	}
	if item.CompletedAt != nil {
		formatted := item.CompletedAt.Format(timeLayout)
		response.CompletedAt = &formatted
	}
	return response
}

// respondData 统一封装成功响应（带数据）。
// 格式：{"data": ...}
func respondData(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

// respondOK 是 respondData 的快捷方式，默认状态码 200 OK。
func respondOK(c *gin.Context, data interface{}) {
	respondData(c, http.StatusOK, data)
}

// respondError 统一封装错误响应。
// 格式：{"error": "message"}
// 这让前端可以统一处理错误逻辑。
func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// respondUnauthorized 统一返回 unauthorized 错误，但在后端日志中记录具体原因。
// 这样可以避免向客户端泄露敏感的错误信息，同时方便后端调试。
// actualReason: 实际的错误原因，会记录到日志中
func respondUnauthorized(c *gin.Context, actualReason string) {
	// 记录具体的错误原因到后端日志
	log.Printf("[Unauthorized] Path: %s, Reason: %s", c.Request.URL.Path, actualReason)
	// 统一返回 unauthorized 给客户端
	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}