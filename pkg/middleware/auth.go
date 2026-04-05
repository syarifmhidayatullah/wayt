package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/project/wayt/pkg/response"
)

// JWTAuth validates Bearer token from Authorization header.
// Used for requests from the admin frontend.
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	secret := []byte(jwtSecret)
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		c.Set("user_id", claims["sub"])
		c.Set("branch_id", claims["branch_id"])
		c.Next()
	}
}

// SuperAdminOnly hanya mengizinkan role superadmin. Harus dipasang setelah JWTAuth.
func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "superadmin" {
			response.Forbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
