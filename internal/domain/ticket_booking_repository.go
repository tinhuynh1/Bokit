package domain

import (
	"context"

	"gorm.io/gorm"
)

type TicketBookingRepository interface {
	AcquireBookingLock(ctx context.Context) error
	ReleaseBookingLock(ctx context.Context) error
	CreateWithTx(ctx context.Context, tx *gorm.DB, booking *TicketBooking) error
	GetByEventID(ctx context.Context, eventID int) ([]TicketBooking, error)
	GetByEmail(ctx context.Context, email string) ([]TicketBooking, error)
	GetByStatus(ctx context.Context, status string) ([]TicketBooking, error)
	CancelBooking(ctx context.Context, bookingIds []int) error
	GetExpiredBooking(ctx context.Context) ([]TicketBooking, error)
}
