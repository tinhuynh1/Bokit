package router

import (
	"booking-svc/internal/handler"
	"booking-svc/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, eventHandler *handler.EventHandler, paymentHandler *handler.PaymentHandler) {
	r.Use(middleware.JWTAuthMiddleware())
	event := r.Group("/events")
	{
		event.GET("", eventHandler.ListEvent)
		event.POST("", eventHandler.CreateEvent)
		event.PATCH("/:event_id", eventHandler.UpdateEvent)
		event.DELETE("/:event_id", eventHandler.DeleteEvent)
		event.POST("/:event_id/booking", eventHandler.BookTicket)
		event.GET("/stats", eventHandler.GetEventStats)
	}

	payment := r.Group("/payment")
	{
		payment.POST("/confirm", paymentHandler.ConfirmPayment)
		payment.POST("/callback", paymentHandler.PaymentCallback)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
