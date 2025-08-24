package dto

type ListEventResponse struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	DateTime         string  `json:"date_time"`
	TicketPrice      float64 `json:"ticket_price"`
	AvailableTickets int     `json:"available_tickets"`
	SoldTickets      int     `json:"sold_tickets"`
}
