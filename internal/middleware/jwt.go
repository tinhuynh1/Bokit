package middleware

import (
	"booking-svc/config"
	"booking-svc/internal/common/response"
	"booking-svc/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware kiểm tra và parse token từ header
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, response.ErrorCodeUnauthorized, "missing_authorization_header"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, response.ErrorCodeUnauthorized, "invalid_token_format"))
			return
		}

		claims, err := utils.VerifyAccessToken(tokenString, config.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, response.ErrorCodeUnauthorized, "invalid_or_expired_token"))
			return
		}

		c.Set("role", claims.Role)
		c.Next()
	}
}
