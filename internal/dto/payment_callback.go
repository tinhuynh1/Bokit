package dto

type PaymentCallbackRequest struct {
	PaymentMethod string `json:"payment_method"`
	Signature     string `json:"signature"`
	TransactionId string `json:"transaction_id"`
	ExtraData     string `json:"extra"`
}
