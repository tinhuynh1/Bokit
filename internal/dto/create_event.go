package dto

import "booking-svc/internal/common/request"

type CreateEventRequest struct {
	request.Request
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	DateTime     string  `json:"date_time" binding:"required"`
	TicketPrice  float64 `json:"ticket_price" binding:"required"`
	TotalTickets int     `json:"total_tickets" binding:"required"`
}
