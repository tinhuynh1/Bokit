package handler

import (
	"booking-svc/internal/common/response"
	"booking-svc/internal/dto"
	eventsvc "booking-svc/internal/service/event"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type EventHandler struct {
	service *eventsvc.EventService
	logger  *zap.Logger
}

func NewEventHandler(svc *eventsvc.EventService, logger *zap.Logger) *EventHandler {
	return &EventHandler{service: svc, logger: logger}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("booking-svc")
	_, span := tracer.Start(ctx, "create_event")
	defer span.End()
	h.logger.Info("create event request", zap.String("trace_id", span.SpanContext().TraceID().String()))
	role, _ := c.Get("role")
	if roleStr, _ := role.(string); roleStr != "admin" {
		c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, response.ErrorCodeForbidden, "access_forbidden"))
		return
	}

	req := dto.CreateEventRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
		return
	}

	err := h.service.CreateEvent(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "create_event_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil, "create_event_success"))
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
		return
	}

	req := dto.UpdateEventRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
		return
	}
	req.ID = id
	if err := h.service.UpdateEvent(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeBadRequest, "update_event_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil, "update_event_success"))
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	ctx := c.Request.Context()
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	h.logger.Info("delete event request", zap.String("trace_id", traceID))
	id, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
		return
	}

	if err := h.service.DeleteEvent(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "delete_event_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil, "delete_event_success"))
}

func (h *EventHandler) ListEvent(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("booking-svc")
	_, span := tracer.Start(ctx, "list_event")
	defer span.End()
	h.logger.Info("list event request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	role, _ := c.Get("role")
	if roleStr, _ := role.(string); roleStr != "admin" {
		c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, response.ErrorCodeForbidden, "access_forbidden"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	from := c.Query("from")
	to := c.Query("to")

	events, paging, err := h.service.ListEvents(ctx, page, pageSize, from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeBadRequest, "list_events_failed"))
		return
	}
	c.JSON(http.StatusOK, response.SuccessWithPaging(events, paging, "list_events_success"))
}

func (h *EventHandler) BookTicket(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
		return
	}

	var req dto.BookingEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
		return
	}
	req.EventID = eventID
	if err := h.service.BookTicket(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeBadRequest, "book_ticket_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil, "book_ticket_success"))
}

func (h *EventHandler) GetEventStats(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("booking-svc")
	_, span := tracer.Start(ctx, "get_event_stats")
	defer span.End()
	h.logger.Info("get event stats request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	role, _ := c.Get("role")
	if roleStr, _ := role.(string); roleStr != "admin" {
		c.JSON(http.StatusForbidden, response.Error(http.StatusForbidden, response.ErrorCodeForbidden, "access_forbidden"))
		return
	}

	from := c.Query("from")
	to := c.Query("to")
	stats, err := h.service.GetEventStats(ctx, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "get_event_stats_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(stats, "get_event_stats_success"))
}
