package bootstrap

import (
	"booking-svc/config"
	"booking-svc/internal/handler"
	"booking-svc/internal/infra/cache"
	"booking-svc/internal/infra/database"
	"booking-svc/internal/job"
	"booking-svc/internal/repository"
	"booking-svc/internal/router"
	"booking-svc/internal/service/event"
	"booking-svc/pkg/logger"
	telemetry "booking-svc/pkg/tracer"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

type App struct {
	engine *gin.Engine
	cfg    *config.Config
}

func NewApp() (*App, error) {
	shutdown := telemetry.InitTracer()
	defer shutdown()
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger.Init()

	db, err := database.NewPostgres(&cfg.Database.Postgres)
	if err != nil {
		logger.L.Error("failed to connect to postgres", zap.Error(err))
		return nil, err
	}

	rdb, err := cache.NewRedis(&cfg.Redis)
	if err != nil {
		logger.L.Error("failed to connect to redis", zap.Error(err))
		return nil, err
	}

	eventRepo := repository.NewEventRepo(db, rdb)
	bookingRepo := repository.NewTicketBookingRepo(db, rdb)
	eventSvc := event.NewEventService(logger.L, eventRepo, bookingRepo)

	eventHandler := handler.NewEventHandler(eventSvc, logger.L)

	cronJob := job.NewCancelBookingJob(eventSvc, logger.L)
	cronJob.Run()

	// // Setup Gin
	engine := gin.New()
	// engine.Use(middleware.TracingMiddleware("booking-svc"))
	// //engine.Use(middleware.LoggingMiddleware()) // custom structured logging
	engine.Use(gin.Recovery())
	engine.Use(otelgin.Middleware("booking-svc"))

	router.SetupRoutes(engine, eventHandler)
	return &App{
		engine: engine,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port)
	return a.engine.Run(addr)
}
