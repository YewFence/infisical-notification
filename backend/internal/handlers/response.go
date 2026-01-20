// 统一响应结构与序列化逻辑。
package handlers

import (
	"net/http"
	"time"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

// TodoResponse is the JSON view of a todo item.
type TodoResponse struct {
	ID          uint       `json:"id"`
	SecretPath  string     `json:"secretPath"`
	IsCompleted bool       `json:"isCompleted"`
	CreatedAt   string     `json:"createdAt"`
	CompletedAt *string    `json:"completedAt"`
}

const timeLayout = time.RFC3339

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

func respondData(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

func respondOK(c *gin.Context, data interface{}) {
	respondData(c, http.StatusOK, data)
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

