package event

import (
	"booking-svc/internal/domain"
	"booking-svc/internal/dto"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock ---
type MockBookingRepo struct {
	mock.Mock
}

func (m *MockBookingRepo) GetBookingById(ctx context.Context, id string) (*domain.TicketBooking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TicketBooking), args.Error(1)
}

func (m *MockBookingRepo) UpdateStatusById(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// --- Tests ---
func TestPaymentCallback(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockBookingRepo)
		svc := NewEventService(nil, nil, nil, nil)

		req := dto.PaymentCallbackRequest{BookingID: 123}
		booking := &domain.TicketBooking{ID: 123}

		mockRepo.On("GetBookingById", ctx, "123").Return(booking, nil)
		mockRepo.On("UpdateStatusById", ctx, "123", string(domain.BookingStatusConfirmed)).Return(nil)

		err := svc.PaymentCallback(ctx, req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get booking error", func(t *testing.T) {
		mockRepo := new(MockBookingRepo)
		svc := NewEventService(nil, nil, nil, nil)

		req := dto.PaymentCallbackRequest{BookingID: 123}
		mockRepo.On("GetBookingById", ctx, 123).Return(nil, errors.New("db error"))

		err := svc.PaymentCallback(ctx, req)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update booking error", func(t *testing.T) {
		mockRepo := new(MockBookingRepo)
		svc := NewEventService(nil, nil, nil, nil)

		req := dto.PaymentCallbackRequest{BookingID: 123}
		booking := &domain.TicketBooking{ID: 123}

		mockRepo.On("GetBookingById", ctx, "123").Return(booking, nil)
		mockRepo.On("UpdateStatusById", ctx, "123", string(domain.BookingStatusConfirmed)).Return(errors.New("update fail"))

		err := svc.PaymentCallback(ctx, req)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
