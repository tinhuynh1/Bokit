package domain

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

type RefundBooking struct {
	EventID    int     `gorm:"column:event_id;not null"`
	TotalPrice float64 `gorm:"column:total_price;not null"`
}

type TicketBooking struct {
	ID         int           `gorm:"primaryKey"`
	EventID    int           `gorm:"column:event_id;not null"`
	Email      string        `gorm:"column:email;not null"`
	Quantity   int           `gorm:"column:quantity;not null"`
	Status     BookingStatus `gorm:"column:status;not null"`
	TotalPrice float64       `gorm:"column:total_price;not null"`
	CreatedAt  time.Time     `gorm:"column:created_at;not null"`
}
