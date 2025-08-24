package payment

import (
	"booking-svc/internal/dto"
	"booking-svc/internal/infra/nats"
	"context"

	"go.uber.org/zap"
)

type PaymentService struct {
	logger *zap.Logger
	nats   *nats.Publisher
}

func NewPaymentService(logger *zap.Logger, nats *nats.Publisher) *PaymentService {
	return &PaymentService{logger: logger, nats: nats}
}

func (s *PaymentService) ConfirmPayment(ctx context.Context, req dto.PaymentConfirmRequest) error {
	//write message to kafka
	s.nats.Publish("payment_confirm", "payment_confirm", req)
	return nil
}

func (s *PaymentService) PaymentCallback(ctx context.Context, req dto.PaymentCallbackRequest) error {
	return nil
}
