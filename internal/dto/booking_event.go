package dto

import "booking-svc/internal/common/request"

type BookingEventRequest struct {
	request.Request
	EventID  int
	Email    string `json:"email" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type BookingEventResponse struct {
	EventID int `json:"event_id"`
}
