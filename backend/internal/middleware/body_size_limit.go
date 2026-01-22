// Package middleware 提供 HTTP 中间件功能。
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodySizeLimit 返回一个限制请求体大小的中间件。
// 如果请求体超过指定的大小限制，将返回 413 状态码（Request Entity Too Large）。
func BodySizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置请求体的最大读取大小
		// 如果请求体超过这个大小，读取操作会失败
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		// 继续处理请求
		c.Next()

		// 检查是否因为请求体过大而产生错误
		if c.Errors.Last() != nil {
			if c.Writer.Status() == http.StatusRequestEntityTooLarge {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "请求体过大",
				})
				c.Abort()
				return
			}
		}
	}
}
