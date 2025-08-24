package dto

type PaymentConfirmRequest struct {
	BookingID     int    `json:"booking_id"`
	PaymentMethod string `json:"payment_method"`
	Signature     string `json:"signature"`
	TransactionID string `json:"transaction_id"`
	ExtraData     string `json:"extra"`
}
