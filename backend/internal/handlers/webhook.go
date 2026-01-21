// Package handlers 包含 Webhook 相关的处理逻辑。
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"backend/internal/repo"
	"backend/internal/signature"

	"github.com/gin-gonic/gin"
)

const (
	// Infisical 发送的签名头
	signatureHeader = "x-infisical-signature"

	// 支持的事件类型
	eventSecretsModified = "secrets.modified"
	eventTest            = "test"
)

// WebhookHandler 专门处理 Webhook 请求。
type WebhookHandler struct {
	repo   *repo.TodoRepository
	secret string // 用于验证签名的密钥
}

// NewWebhookHandler 创建 WebhookHandler 实例。
func NewWebhookHandler(repo *repo.TodoRepository, secret string) *WebhookHandler {
	return &WebhookHandler{repo: repo, secret: strings.TrimSpace(secret)}
}

// webhookPayload 定义了 Infisical Webhook 的 JSON 载荷结构。
type webhookPayload struct {
	Event   string `json:"event"`
	Project struct {
		SecretPath string `json:"secretPath"`
		// 以下字段目前未使用，但保留方便将来扩展
		ProjectID    string `json:"projectId"`
		ProjectName  string `json:"projectName"`
		Environment  string `json:"environment"`
		SecretName   string `json:"secretName"`
		ReminderNote string `json:"reminderNote"`
	} `json:"project"`
	Timestamp int64 `json:"timestamp"`
}

// Handle 处理 Webhook 请求的主要逻辑。
//
//	@Summary		接收 Infisical Webhook
//	@Description	接收来自 Infisical 的 Webhook 通知并创建或更新待办事项
//	@Tags			webhook
//	@Accept			json
//	@Produce		json
//	@Param			X-Infisical-Signature	header		string					true	"Webhook 签名"
//	@Param			payload					body		webhookPayload			true	"Webhook 载荷"
//	@Success		200						{object}	map[string]interface{}	"成功处理 Webhook"
//	@Failure		400						{object}	map[string]string		"请求参数错误"
//	@Failure		401						{object}	map[string]string		"签名验证失败"
//	@Failure		500						{object}	map[string]string		"服务器内部错误"
//	@Router			/webhook [post]
//	@Security		WebhookSecret
func (h *WebhookHandler) Handle(c *gin.Context) {
	// 1. 获取原始请求体 (Raw Data)
	// 验证签名需要原始的字节流，而不是解析后的 JSON 对象。
	// 任何对 JSON 的微小改动（如空格）都会导致签名验证失败。
	bodyBytes, err := c.GetRawData()
	if err != nil {
		respondError(c, http.StatusBadRequest, "read body failed")
		return
	}

	bodyText := strings.TrimSpace(string(bodyBytes))
	if bodyText == "" {
		respondError(c, http.StatusBadRequest, "empty body")
		return
	}

	// 2. 检查系统是否配置了 Webhook Secret
	if strings.TrimSpace(h.secret) == "" {
		respondError(c, http.StatusInternalServerError, "missing webhook secret")
		return
	}

	// 3. 获取签名头
	signatureHeaderValue := strings.TrimSpace(c.GetHeader(signatureHeader))
	if signatureHeaderValue == "" {
		respondError(c, http.StatusBadRequest, "missing signature header")
		return
	}

	// 4. 验证签名
	// 调用 signature 包的逻辑，确保请求确实来自 Infisical 且未被篡改。
	if err := signature.VerifySignature(bodyText, signatureHeaderValue, h.secret, time.Now().UTC()); err != nil {
		respondError(c, http.StatusUnauthorized, "invalid signature")
		return
	}

	// 5. 解析 JSON 载荷
	var payload webhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	// 6. 过滤事件类型
	if !isSupportedEvent(payload.Event) {
		respondOK(c, "ignored")
		return
	}

	// 测试事件直接返回 OK
	if payload.Event == eventTest {
		respondOK(c, "ok")
		return
	}

	// 7. 处理业务逻辑
	// 提取 secretPath，如果没有则默认为根路径 "/"
	secretPath := strings.TrimSpace(payload.Project.SecretPath)
	if secretPath == "" {
		secretPath = "/"
	}

	// 更新或插入 Todo 项
	item, err := h.repo.UpsertFromWebhook(secretPath, time.Now().UTC())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "upsert todo failed")
		return
	}

	respondOK(c, toTodoResponse(item))
}

func isSupportedEvent(event string) bool {
	switch event {
	case eventSecretsModified, eventTest:
		return true
	default:
		return false
	}
}
