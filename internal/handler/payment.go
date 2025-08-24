package handler

import (
	"booking-svc/internal/common/response"
	"booking-svc/internal/dto"
	eventsvc "booking-svc/internal/service/event"
	paymentsvc "booking-svc/internal/service/payment"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	service    *eventsvc.EventService
	paymentSvc *paymentsvc.PaymentService
	logger     *zap.Logger
}

func NewPaymentHandler(svc *eventsvc.EventService, paymentSvc *paymentsvc.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{service: svc, paymentSvc: paymentSvc, logger: logger}
}

func (h *PaymentHandler) ConfirmPayment(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("booking-svc")
	_, span := tracer.Start(ctx, "confirm_payment")
	defer span.End()
	h.logger.Info("confirm payment request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	req := dto.PaymentConfirmRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
		return
	}

	if err := h.paymentSvc.ConfirmPayment(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "confirm_payment_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil, "confirm_payment_success"))
}

func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("booking-svc")
	_, span := tracer.Start(ctx, "payment_callback")
	defer span.End()
	h.logger.Info("payment callback request", zap.String("trace_id", span.SpanContext().TraceID().String()))

	req := dto.PaymentCallbackRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
		return
	}

	if err := h.service.PaymentCallback(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "payment_callback_failed"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil, "payment_callback_success"))
}
