package dto

type UpdateEventRequest struct {
	ID               int
	Name             string  `json:"name,omitempty"`
	Description      string  `json:"description,omitempty"`
	DateTime         string  `json:"date_time,omitempty"`
	TicketPrice      float64 `json:"ticket_price,omitempty"`
	AvailableTickets *int    `json:"available_tickets,omitempty"`
}
