package middleware

import (
	"booking-svc/internal/common/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware giúp không crash server nếu có panic
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Panic: %v", rec)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Status:     "error",
					Code:       http.StatusInternalServerError,
					MessageKey: "internal_server_error",
					ErrorCode:  response.ErrorCodeInternalServer,
				})
			}
		}()
		c.Next()
	}
}
