package handler

import (
	"net/http"
	"quiz-svc/internal/common/response"
	"quiz-svc/internal/service"

	"quiz-svc/internal/dto"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type QuizHandler struct {
	service *service.QuizService
	logger  *zap.Logger
}

func NewQuizHandler(service *service.QuizService, logger *zap.Logger) *QuizHandler {
	return &QuizHandler{service: service, logger: logger}
}

func (h *QuizHandler) ListQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("quiz-svc")
	_, span := tracer.Start(ctx, "list_quiz")
	defer span.End()
	h.logger.Info("list quiz request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	quizzes, err := h.service.ListQuiz(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "list_quiz_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(quizzes, "list_quiz_success"))
}

func (h *QuizHandler) GetQuizById(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("quiz-svc")
	_, span := tracer.Start(ctx, "get_quiz_by_id")
	defer span.End()
	h.logger.Info("get quiz by id request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	quiz, err := h.service.GetQuizById(ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "get_quiz_by_id_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(quiz, "get_quiz_by_id_success"))
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("quiz-svc")
	_, span := tracer.Start(ctx, "update_quiz")
	defer span.End()
	h.logger.Info("update quiz request", zap.String("trace_id", span.SpanContext().TraceID().String()))
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("quiz-svc")
	_, span := tracer.Start(ctx, "create_quiz")
	defer span.End()
	h.logger.Info("create quiz request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	quiz := &dto.CreateQuizRequest{}
	if err := c.ShouldBindJSON(quiz); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeBadRequest, "invalid_request"))
		return
	}

	err := h.service.CreateQuiz(ctx, quiz)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "create_quiz_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil, "create_quiz_success"))
}
