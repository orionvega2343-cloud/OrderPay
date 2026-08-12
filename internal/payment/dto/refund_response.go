package dto

import "time"

type RefundResponse struct {
	Id        int       `json:"id"`
	PaymentId int       `json:"payment_id"`
	Amount    int       `json:"amount"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
