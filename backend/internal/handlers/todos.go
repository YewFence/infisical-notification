// TODO CRUD 处理器，提供列表与状态更新。
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

// TodoHandler handles CRUD requests for todo items.
type TodoHandler struct {
	repo *repo.TodoRepository
}

// NewTodoHandler builds a new handler instance.
func NewTodoHandler(repo *repo.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

type todoInput struct {
	SecretPath string `json:"secretPath"`
}

func (h *TodoHandler) List(c *gin.Context) {
	items, err := h.repo.List()
	if err != nil {
		respondError(c, http.StatusInternalServerError, "list todos failed")
		return
	}

	response := make([]TodoResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toTodoResponse(item))
	}
	respondOK(c, response)
}

func (h *TodoHandler) Create(c *gin.Context) {
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

	item, err := h.repo.Create(secretPath, time.Now().UTC())
	if err != nil {
		if isUniqueConstraintError(err) {
			respondError(c, http.StatusConflict, "secretPath already exists")
			return
		}
		respondError(c, http.StatusInternalServerError, "create todo failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

func (h *TodoHandler) Update(c *gin.Context) {
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
		if isUniqueConstraintError(err) {
			respondError(c, http.StatusConflict, "secretPath already exists")
			return
		}
		respondError(c, http.StatusInternalServerError, "update todo failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

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

func parseID(c *gin.Context) (uint, bool) {
	idText := strings.TrimSpace(c.Param("id"))
	idValue, err := strconv.ParseUint(idText, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(idValue), true
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "unique")
}

