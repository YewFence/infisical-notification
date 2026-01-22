// Package middleware 提供 HTTP 中间件功能。
package middleware

import (
	"net/http"

	"backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

// BodySizeLimit 返回一个限制请求体大小的中间件。
// 如果请求体超过指定的大小限制，将返回 413 状态码（Request Entity Too Large）。
func BodySizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先检查 Content-Length 头，快速拒绝明显过大的请求
		if c.Request.ContentLength > maxSize {
			handlers.RespondError(c, http.StatusRequestEntityTooLarge, "请求体过大")
			c.Abort()
			return
		}

		// 包装请求体，限制实际读取的大小
		// 这可以防止客户端发送虚假的 Content-Length 头
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()
	}
}
