package router

import (
	"quiz-svc/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine,
	quizHandler *handler.QuizHandler,
	sessionHandler *handler.SessionHandler,
	wsHandler *handler.WSHandler) {
	//r.Use(middleware.JWTAuthMiddleware())
	v1 := r.Group("api/v1")

	// Quiz endpoints
	quiz := v1.Group("quizzes")
	{
		quiz.GET("", quizHandler.ListQuiz)
		quiz.GET("/:id", quizHandler.GetQuizById)
		quiz.PUT("/:id", quizHandler.UpdateQuiz)
		quiz.POST("", quizHandler.CreateQuiz) //OK
	}

	// Session endpoints
	session := v1.Group("sessions")
	{
		//session.GET("/:code", sessionHandler.GetSession)
		session.POST("", sessionHandler.CreateSession) //OK
		//session.POST("/:code/start", sessionHandler.StartSession)
		//session.GET("/:code/leaderboard", sessionHandler.GetLeaderboard)
	}

	// WebSocket endpoint
	r.GET("/ws", wsHandler.Ws)

	// Serve 1 file
	r.StaticFile("/", "./test_client.html")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
