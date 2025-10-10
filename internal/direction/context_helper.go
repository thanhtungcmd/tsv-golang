package direction

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		ctx := context.WithValue(c.Request.Context(), "token", tokenString)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
