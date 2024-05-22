package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 函数是一个 Gin 中间件，用于处理跨域请求。
// 它会在每个请求之前添加必要的响应头，以允许跨域访问。
// 具体而言，它会设置以下响应头：
// - Access-Control-Allow-Origin: 允许所有来源的请求访问
// - Access-Control-Allow-Headers: 允许的请求头字段
// - Access-Control-Allow-Methods: 允许的请求方法
// - Access-Control-Expose-Headers: 允许客户端访问的响应头字段
// - Access-Control-Allow-Credentials: 允许发送身份凭证（如 cookies）
// 如果请求方法为 OPTIONS，则会中断请求并返回状态码 204。
// 否则，会继续处理请求。

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
