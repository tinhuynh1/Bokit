package payment

import (
	"booking-svc/internal/domain"
	"booking-svc/internal/dto"
	"booking-svc/internal/infra/message_broker"
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type PaymentService struct {
	logger      *zap.Logger
	bookingRepo domain.TicketBookingRepository
	nats        *message_broker.Publisher
	consumer    *message_broker.Consumer
}

func NewPaymentService(
	logger *zap.Logger,
	bookingRepo domain.TicketBookingRepository,
	nats *message_broker.Publisher,
	consumer *message_broker.Consumer,
) *PaymentService {
	return &PaymentService{
		logger:      logger,
		bookingRepo: bookingRepo,
		nats:        nats,
		consumer:    consumer,
	}
}

func (s *PaymentService) ConfirmPayment(ctx context.Context, req dto.PaymentConfirmRequest) error {
	exists, err := s.bookingRepo.IsExistsBookingProcessing(ctx, req.BookingID)
	if err != nil {
		return fmt.Errorf("failed to check idempotence: %w", err)
	}
	if exists {
		s.logger.Info("booking is already processing", zap.Int("booking_id", req.BookingID))
		return nil
	}
	booking, err := s.bookingRepo.GetBookingById(ctx, req.BookingID)
	if err != nil {
		s.logger.Error("get booking by id failed", zap.Error(err))
		return err
	}
	if booking.Status != domain.BookingStatusPending {
		s.logger.Error("booking status is not pending", zap.Int("booking_id", booking.ID))
		return errors.New("booking status is not pending")
	}
	//write message to nats
	msg := message_broker.Message{
		Type: "payment_confirm",
		Data: req,
		// Timestamp: time.Now(),
		// Source:    "booking-service",
	}
	err = s.nats.Publish("payment_confirm", msg.Type, msg.Data)
	if err != nil {
		s.logger.Error("failed to publish message", zap.Error(err))
		return err
	}
	return nil
}

func (s *PaymentService) PaymentCallback(ctx context.Context, req dto.PaymentCallbackRequest) error {
	return nil
}

func (s *PaymentService) StartPaymentConsumer() error {

	err := s.consumer.Subscribe("payment_confirm", func(msg *nats.Msg) {
		bookingID := string(msg.Data)
		s.logger.Info("Simulate payment processing", zap.String("booking_id", bookingID))

		// ✅ Xử lý message
		msg.Ack()
	})

	if err != nil {
		s.logger.Error("failed to subscribe to payment_confirm", zap.Error(err))
		return err
	}

	return nil
}
