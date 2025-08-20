package domain

import "time"

type Event struct {
	ID           int        `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string     `gorm:"column:name"`
	Description  string     `gorm:"column:description"`
	DateTime     time.Time  `gorm:"column:date_time"`
	TicketPrice  float64    `gorm:"column:ticket_price"`
	TotalTickets int        `gorm:"column:total_tickets"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}
