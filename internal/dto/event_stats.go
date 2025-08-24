package dto

type EventStatsResponse struct {
	Period           string       `json:"period"`
	From             string       `json:"from"`
	To               string       `json:"to"`
	TotalEvents      int          `json:"total_events"`
	TotalTicketsSold int          `json:"total_tickets_sold"`
	TotalRevenue     float64      `json:"total_revenue"`
	Events           []EventStats `json:"events"`
}

type EventStats struct {
	EventID           int     `json:"event_id"`
	EventName         string  `json:"event_name"`
	DateTime          string  `json:"date_time"`
	TotalTickets      int     `json:"total_tickets"`
	TotalTicketsSold  int     `json:"total_tickets_sold"`
	AvailableTickets  int     `json:"available_tickets"`
	EstimatedRevenue  float64 `json:"estimated_revenue"`
	TotalBookings     int     `json:"total_bookings"`
	ConfirmedBookings int     `json:"confirmed_bookings"`
	PendingBookings   int     `json:"pending_bookings"`
	CancelledBookings int     `json:"cancelled_bookings"`
}
