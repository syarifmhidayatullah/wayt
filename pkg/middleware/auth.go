package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/project/wayt/pkg/response"
)

func InternalAuth(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Internal-Key")
		if key == "" || key != apiKey {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
