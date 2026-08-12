package domain

import "time"

type Refund struct {
	Id        int       `db:"id"`
	PaymentId int       `db:"payment_id"`
	Amount    int       `db:"amount"`
	Reason    string    `db:"reason"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}
