package bootstrap

import (
	"context"
	"fmt"
	"quiz-svc/config"
	"quiz-svc/internal/handler"
	"quiz-svc/internal/infra/cache"
	"quiz-svc/internal/infra/database"
	"quiz-svc/internal/infra/message_broker"
	"quiz-svc/internal/repository"
	"quiz-svc/internal/router"
	"quiz-svc/internal/service"
	"quiz-svc/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	engine *gin.Engine
	cfg    *config.Config
}

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	//telemetry.InitTracer()

	logger.Init()

	db, err := database.NewMongoDB(&cfg.Database.MongoDB)
	if err != nil {
		logger.L.Error("failed to connect to mongodb", zap.Error(err))
		return nil, err
	}

	quizRepo := repository.NewQuizCollection(db)
	sessionRepo := repository.NewSessionCollection(db)

	redisClient, err := cache.NewRedis(&cfg.Redis)
	if err != nil {
		logger.L.Error("failed to connect to redis", zap.Error(err))
		return nil, err
	}

	// Initialize Redis repository
	redisRepo := repository.NewRedisRepo(redisClient)

	// Initialize services

	quizSvc := service.NewQuizService(logger.L, quizRepo)
	sessionSvc := service.NewSessionService(sessionRepo, redisRepo)
	netSvc := service.NewConnectionService(sessionSvc)

	// Initialize handlers
	quizHandler := handler.NewQuizHandler(quizSvc, logger.L)
	sessionHandler := handler.NewSessionHandler(sessionSvc, quizSvc, logger.L)
	wsHandler := handler.NewWSHandler(netSvc, logger.L)

	// Setup Gin
	engine := gin.New()
	engine.Use(gin.Recovery())
	//	engine.Use(otelgin.Middleware("quiz-svc"))

	router.SetupRoutes(engine, quizHandler, sessionHandler, wsHandler)

	go func() {
		for {
			fmt.Println("HandleParticipantJoined")
			sessionSvc.HandleParticipantJoined(context.Background())
			fmt.Println("HandleParticipantJoined done")
		}
	}()
	return &App{
		engine: engine,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port)
	return a.engine.Run(addr)
}

func (a *App) Cleanup() {
	database.Close()
	cache.Close()
	message_broker.ClosePublisher()
	message_broker.CloseConsumer()
}
