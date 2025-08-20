package handler

import (
	"booking-svc/internal/common/response"
	eventsvc "booking-svc/internal/service/event"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	service *eventsvc.EventService
	logger  *zap.Logger
}

func NewPaymentHandler(svc *eventsvc.EventService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{service: svc, logger: logger}
}

func (h *PaymentHandler) CallbackPayment(c *gin.Context) {
	// ctx := c.Request.Context()
	// req := dto.PaymentCallbackRequest{}
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, response.ErrorCodeInvalidRequest, "invalid_request"))
	// 	return
	// }

	// err := h.service.CallbackPayment(ctx, req)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, response.ErrorCodeInternalServer, "create_event_failed"))
	// 	return
	// }
	c.JSON(http.StatusOK, response.Success(nil, "create_event_success"))
}
