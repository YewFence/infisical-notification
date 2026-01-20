// Infisical webhook handler that verifies signatures and forwards notifications to Apprise.
package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

const (
	signatureHeader      = "x-infisical-signature"
	eventSecretsModified = "secrets.modified"
	eventTest            = "test"
	replayWindow         = 5 * time.Minute
)

// messageTemplate is embedded at build time to avoid runtime file dependencies.
//
//go:embed message.md
var messageTemplate string

type webhookPayload struct {
	Event   string `json:"event"`
	Project struct {
		SecretPath string `json:"secretPath"`
	} `json:"project"`
}

type messageData struct {
	Event      string
	SecretPath string
}

func Main(context openruntimes.Context) openruntimes.Response {
	bodyText := context.Req.BodyText()
	if strings.TrimSpace(bodyText) == "" {
		context.Log("empty request body")
		return context.Res.Text("empty body")
	}

	secret := strings.TrimSpace(os.Getenv("INFISICAL_WEBHOOK_SECRET"))
	if secret == "" {
		context.Log("missing INFISICAL_WEBHOOK_SECRET")
		return context.Res.Text("missing secret")
	}

	signature := strings.TrimSpace(context.Req.Headers[signatureHeader])
	if signature == "" {
		context.Log("missing x-infisical-signature header")
		return context.Res.Text("missing signature")
	}

	if err := verifySignature(bodyText, signature, secret, time.Now()); err != nil {
		context.Log(fmt.Sprintf("signature verification failed: %v", err))
		return context.Res.Text("invalid signature")
	}

	var payload webhookPayload
	if err := json.Unmarshal([]byte(bodyText), &payload); err != nil {
		context.Log(fmt.Sprintf("invalid json payload: %v", err))
		return context.Res.Text("invalid payload")
	}

	if !isSupportedEvent(payload.Event) {
		return context.Res.Text("ignored")
	}

	appriseURL := strings.TrimSpace(os.Getenv("APPRISE_URL"))
	if appriseURL == "" {
		context.Log("missing APPRISE_URL")
		return context.Res.Text("missing apprise url")
	}

	notificationURLs := strings.TrimSpace(os.Getenv("NOTIFICATION_URLS"))
	if notificationURLs == "" {
		context.Log("missing NOTIFICATION_URLS")
		return context.Res.Text("missing notification urls")
	}

	secretPath := strings.TrimSpace(payload.Project.SecretPath)
	if secretPath == "" {
		secretPath = "/"
	}

	messageBody, err := renderMessage(payload.Event, secretPath)
	if err != nil {
		context.Log(fmt.Sprintf("render message failed: %v", err))
		return context.Res.Text("render error")
	}

	title := "Infisical secrets updated"
	if payload.Event == eventTest {
		title = "Infisical webhook test"
	}
	if err := sendApprise(appriseURL, notificationURLs, title, messageBody); err != nil {
		context.Log(fmt.Sprintf("apprise request failed: %v", err))
		return context.Res.Text("notify failed")
	}

	return context.Res.Text("ok")
}

func verifySignature(bodyText, headerValue, secret string, now time.Time) error {
	timestamp, signature, err := parseSignatureHeader(headerValue)
	if err != nil {
		return err
	}

	if !isTimestampFresh(timestamp, now) {
		return errors.New("timestamp out of range")
	}

	expected := computeHMAC(bodyText, secret)
	decoded, err := decodeSignature(signature)
	if err != nil {
		return err
	}

	if !hmac.Equal(decoded, expected) {
		return errors.New("signature mismatch")
	}

	return nil
}

func parseSignatureHeader(headerValue string) (int64, string, error) {
	trimmed := strings.TrimSpace(headerValue)
	if trimmed == "" {
		return 0, "", errors.New("invalid signature header format")
	}

	parts := strings.SplitN(trimmed, ";", 2)
	if len(parts) != 2 {
		return 0, "", errors.New("invalid signature header format")
	}

	timestampPart := strings.TrimSpace(parts[0])
	signature := strings.TrimSpace(parts[1])
	if !strings.HasPrefix(timestampPart, "t=") || signature == "" {
		return 0, "", errors.New("invalid signature header format")
	}

	timestampText := strings.TrimPrefix(timestampPart, "t=")
	if timestampText == "" {
		return 0, "", errors.New("invalid signature header format")
	}

	timestamp, err := strconv.ParseInt(timestampText, 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid timestamp: %w", err)
	}

	if timestamp > 1_000_000_000_000 {
		timestamp = timestamp / 1000
	}

	return timestamp, signature, nil
}

func isTimestampFresh(timestamp int64, now time.Time) bool {
	signedAt := time.Unix(timestamp, 0)
	if signedAt.IsZero() {
		return false
	}

	diff := now.Sub(signedAt)
	if diff < 0 {
		diff = -diff
	}

	return diff <= replayWindow
}

func computeHMAC(payload, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return mac.Sum(nil)
}

func decodeSignature(signature string) ([]byte, error) {
	trimmed := strings.TrimSpace(signature)
	trimmed = strings.TrimPrefix(trimmed, "sha256=")

	if decoded, err := hex.DecodeString(trimmed); err == nil {
		return decoded, nil
	}

	if decoded, err := base64.StdEncoding.DecodeString(trimmed); err == nil {
		return decoded, nil
	}

	return nil, errors.New("unsupported signature encoding")
}

func renderMessage(event, secretPath string) (string, error) {
	tmpl, err := template.New("message").Option("missingkey=error").Parse(messageTemplate)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if err := tmpl.Execute(&builder, messageData{Event: event, SecretPath: secretPath}); err != nil {
		return "", err
	}

	return builder.String(), nil
}

func isSupportedEvent(event string) bool {
	switch event {
	case eventSecretsModified, eventTest:
		return true
	default:
		return false
	}
}

func sendApprise(appriseURL, notificationURLs, title, body string) error {
	payload := map[string]string{
		"urls":  notificationURLs,
		"body":  body,
		"title": title,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, appriseURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		bodyText, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("apprise status %d: %s", resp.StatusCode, strings.TrimSpace(string(bodyText)))
	}

	return nil
}
