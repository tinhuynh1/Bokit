package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func TracingMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.FullPath()) // tạo span cho request
		defer span.End()

		// lưu lại context mới vào request
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
