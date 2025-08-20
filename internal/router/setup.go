package router

import (
	"booking-svc/internal/handler"
	"booking-svc/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, eventHandler *handler.EventHandler) {
	r.Use(middleware.JWTAuthMiddleware())
	event := r.Group("/events")
	{
		event.GET("", eventHandler.ListEvent)
		event.POST("", eventHandler.CreateEvent)
		event.PATCH("/:event_id", eventHandler.UpdateEvent)
		event.DELETE("/:event_id", eventHandler.DeleteEvent)
		event.POST("/:event_id/booking", eventHandler.BookTicket)
	}

	// payment := r.Group("/payment")
	// {
	// 	payment.POST("/callback", eventHandler.CallbackPayment)
	// }

	// stats := r.Group("/stats")
	// {
	// 	stats.GET("/event/:event_id", eventHandler.GetEventStats)
	// 	stats.GET("/revenue", eventHandler.GetRevenueStats)
	// }

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
