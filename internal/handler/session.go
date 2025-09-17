package handler

import (
	"net/http"
	"quiz-svc/internal/common/response"
	"quiz-svc/internal/service"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type SessionHandler struct {
	sessionService *service.SessionService
	quizService    *service.QuizService
	logger         *zap.Logger
}

func NewSessionHandler(sessionService *service.SessionService, quizService *service.QuizService, logger *zap.Logger) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		quizService:    quizService,
		logger:         logger,
	}
}

type CreateSessionRequest struct {
	QuizID string `json:"quiz_id" binding:"required"`
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("quiz-svc")
	_, span := tracer.Start(ctx, "create_session")
	defer span.End()
	h.logger.Info("create session request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeBadRequest, "invalid_request"))
		return
	}

	// Get quiz by ID
	quiz, err := h.quizService.GetQuizById(ctx, req.QuizID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, response.ErrorCodeNotFound, "quiz_not_found"))
		return
	}

	// Create session with quiz
	session, err := h.sessionService.CreateSession(ctx, quiz)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "create_session_failed"))
		return
	}

	c.JSON(http.StatusOK, response.Success(session, "create_session_success"))
}

func (h *SessionHandler) GetSession(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("quiz-svc")
	_, span := tracer.Start(ctx, "get_session")
	defer span.End()
	h.logger.Info("get session request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	code := c.Param("code")
	session, err := h.sessionService.GetSession(ctx, code)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, response.ErrorCodeNotFound, "session_not_found"))
		return
	}
	c.JSON(http.StatusOK, response.Success(session, "get_session_success"))
}

type StartSessionRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *SessionHandler) StartSession(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("quiz-svc")
	_, span := tracer.Start(ctx, "start_session")
	defer span.End()
	h.logger.Info("start session request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	code := c.Param("code")

	err := h.sessionService.StartSession(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil, "start_session_success"))
}

func (h *SessionHandler) GetLeaderboard(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("quiz-svc")
	_, span := tracer.Start(ctx, "get_leaderboard")
	defer span.End()
	h.logger.Info("get leaderboard request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	code := c.Param("code")
	leaderboard, err := h.sessionService.GetLeaderboard(ctx, code)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(http.StatusNotFound, response.ErrorCodeNotFound, "session_not_found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(leaderboard, "get_leaderboard_success"))
}
