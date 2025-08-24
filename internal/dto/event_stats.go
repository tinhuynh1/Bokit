package dto

type EventStatsResponse struct {
	TotalEvents      int          `json:"total_events"`
	TotalTicketsSold int          `json:"total_tickets_sold"`
	TotalRevenue     float64      `json:"total_revenue"`
	Events           []EventStats `json:"events"`
}

type EventStats struct {
	EventID          int     `json:"event_id"`
	EventName        string  `json:"event_name"`
	TotalTicketsSold int     `json:"total_tickets_sold"`
	TotalRevenue     float64 `json:"total_revenue"`
}
