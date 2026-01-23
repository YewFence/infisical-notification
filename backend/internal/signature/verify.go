// Package signature 提供了 HMAC 签名验证功能。
// 这里的逻辑复用了 Infisical 官方文档推荐的验证算法。
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

// replayWindow 定义了时间戳的有效窗口（5分钟）。
// 用于防止重放攻击 (Replay Attack)。
const replayWindow = 5 * time.Minute

// VerifySignature 验证 Infisical 的 Webhook 签名。
// 参数：
// - bodyText: HTTP 请求体的原始内容。
// - headerValue: x-infisical-signature 头的值。
// - secret: 配置的 Webhook Secret。
// - now: 当前时间，用于校验时间戳。
func VerifySignature(bodyText, headerValue, secret string, now time.Time) error {
	// 1. 解析签名头
	// 格式通常为: t=1234567890;sha256=abcdef...
	timestamp, signature, err := parseSignatureHeader(headerValue)
	if err != nil {
		return err
	}

	// 2. 校验时间戳新鲜度
	// 确保请求是在最近 5 分钟内发出的。
	if !isTimestampFresh(timestamp, now) {
		return errors.New("timestamp out of range")
	}

	// 3. 计算期望的 HMAC 值
	// 使用本地 Secret 对 "timestamp.payload" 进行 HMAC-SHA256 哈希运算。
	expected := computeHMAC(timestamp, bodyText, secret)

	// 4. 解码请求中的签名
	// 签名可能是 Hex 或 Base64 编码。
	decoded, err := decodeSignature(signature)
	if err != nil {
		return err
	}

	// 5. 比较签名
	// 使用 hmac.Equal 防止时序攻击 (Timing Attack)。
	if !hmac.Equal(decoded, expected) {
		return errors.New("signature mismatch")
	}

	return nil
}

// parseSignatureHeader 解析类似 "t=123456;sha256=xyz" 的头。
func parseSignatureHeader(headerValue string) (int64, string, error) {
	trimmed := strings.TrimSpace(headerValue)
	if trimmed == "" {
		return 0, "", errors.New("invalid signature header format")
	}

	// 分割时间戳和签名部分
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

	// 简单的逻辑判断时间戳单位（毫秒 vs 秒），兼容不同格式
	if timestamp > 1_000_000_000_000 {
		timestamp = timestamp / 1000
	}

	return timestamp, signature, nil
}

// isTimestampFresh 检查时间戳是否在允许的误差范围内。
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

// computeHMAC 使用 SHA256 计算 HMAC。
// 按照 Infisical 规范，签名对象是 "<timestamp>.<payload>" 格式的字符串。
func computeHMAC(timestamp int64, payload, secret string) []byte {
	// 构造签名字符串: timestamp.payload
	signedPayload := fmt.Sprintf("%d.%s", timestamp, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(signedPayload))
	return mac.Sum(nil)
}

// decodeSignature 尝试解码 Hex 或 Base64 格式的签名。
func decodeSignature(signature string) ([]byte, error) {
	trimmed := strings.TrimSpace(signature)
	// 移除可能的前缀
	trimmed = strings.TrimPrefix(trimmed, "sha256=")

	// 尝试 Base64 解码
	if decoded, err := base64.StdEncoding.DecodeString(trimmed); err == nil {
		return decoded, nil
	}

	// 尝试 Hex 解码
	if decoded, err := hex.DecodeString(trimmed); err == nil {
		return decoded, nil
	}

	return nil, errors.New("unsupported signature encoding")
}
