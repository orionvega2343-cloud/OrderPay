package dto

type UpdatePaymentRequest struct {
	Status string `json:"status"`
	Method string `json:"method"`
}
