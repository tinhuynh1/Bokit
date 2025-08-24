package event

import (
	"booking-svc/internal/domain"
	"booking-svc/internal/infra/message_broker"
)

type Publisher struct {
	nats *message_broker.Publisher
}

func NewPublisher(nats *message_broker.Publisher) *Publisher {
	return &Publisher{nats: nats}
}

func (p *Publisher) PublishBookingCreated(booking *domain.TicketBooking) error {
	data := map[string]interface{}{
		"booking_id":  booking.ID,
		"event_id":    booking.EventID,
		"email":       booking.Email,
		"quantity":    booking.Quantity,
		"total_price": booking.TotalPrice,
		"status":      booking.Status,
	}

	return p.nats.Publish("booking.events", "booking.created", data)
}

func (p *Publisher) PublishPaymentReceived(bookingID int, paymentMethod string) error {
	data := map[string]interface{}{
		"booking_id":     bookingID,
		"payment_method": paymentMethod,
	}

	return p.nats.Publish("payment.events", "payment.received", data)
}
