// 签名校验逻辑，复用 Infisical 规则。
package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const replayWindow = 5 * time.Minute

// VerifySignature validates the Infisical signature header.
func VerifySignature(bodyText, headerValue, secret string, now time.Time) error {
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

