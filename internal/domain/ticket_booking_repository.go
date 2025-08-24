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
	UpdateStatusByIds(ctx context.Context, tx *gorm.DB, bookingIds []int, status string) error
	UpdateStatusById(ctx context.Context, bookingID int, status string) error
	GetExpiredBooking(ctx context.Context, tx *gorm.DB) ([]TicketBooking, error)
	GetBookingsByEventIds(ctx context.Context, eventIds []int) ([]TicketBooking, error)
	GetBookingById(ctx context.Context, bookingID int) (TicketBooking, error)
}
