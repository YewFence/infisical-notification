// Webhook 处理器，负责校验签名并入库。
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
	signatureHeader      = "x-infisical-signature"
	eventSecretsModified = "secrets.modified"
	eventTest            = "test"
)

// WebhookHandler handles Infisical webhook events.
type WebhookHandler struct {
	repo   *repo.TodoRepository
	secret string
}

// NewWebhookHandler constructs a webhook handler.
func NewWebhookHandler(repo *repo.TodoRepository, secret string) *WebhookHandler {
	return &WebhookHandler{repo: repo, secret: strings.TrimSpace(secret)}
}

type webhookPayload struct {
	Event   string `json:"event"`
	Project struct {
		SecretPath string `json:"secretPath"`
	} `json:"project"`
	Timestamp int64 `json:"timestamp"`
}

func (h *WebhookHandler) Handle(c *gin.Context) {
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

	if strings.TrimSpace(h.secret) == "" {
		respondError(c, http.StatusInternalServerError, "missing webhook secret")
		return
	}

	signatureHeaderValue := strings.TrimSpace(c.GetHeader(signatureHeader))
	if signatureHeaderValue == "" {
		respondError(c, http.StatusBadRequest, "missing signature header")
		return
	}

	if err := signature.VerifySignature(bodyText, signatureHeaderValue, h.secret, time.Now().UTC()); err != nil {
		respondError(c, http.StatusUnauthorized, "invalid signature")
		return
	}

	var payload webhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	if !isSupportedEvent(payload.Event) {
		respondOK(c, "ignored")
		return
	}

	if payload.Event == eventTest {
		respondOK(c, "ok")
		return
	}

	secretPath := strings.TrimSpace(payload.Project.SecretPath)
	if secretPath == "" {
		secretPath = "/"
	}

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

