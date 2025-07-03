package middlewares

import (
	"net/http"
	"strings"

	"POS-BE/libraries/helpers/utils/authorizer"

	"github.com/gin-gonic/gin"
)

func Authorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		claims, err := authorizer.ValidateToken(fields[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Set claims like user_id into Gin context
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}

		c.Next()
	}
}
