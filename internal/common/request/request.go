package request

type Request struct {
	RequestID string `json:"request_id"`
	Timestamp int64  `json:"timestamp"`
}
