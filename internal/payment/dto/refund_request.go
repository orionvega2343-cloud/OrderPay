package dto

type RefundRequest struct {
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}
