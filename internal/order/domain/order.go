package domain

import (
	"OrderPay/internal/order/domain/errs"
	"time"
)

type Order struct {
	Id          int       `db:"id"`
	UserId      string    `db:"user_id"`
	Status      string    `db:"status"`
	TotalAmount int       `db:"total_amount"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (o *Order) TransitionStatus(next string) error {
	switch o.Status {
	case "created":
		switch next {
		case "cancelled", "paid":
			o.Status = next
		default:
			return errs.ErrInvalidTransition
		}
	case "paid":
		switch next {
		case "completed", "cancelled":
			o.Status = next
		default:
			return errs.ErrInvalidTransition
		}
	case "cancelled", "completed":
		return errs.ErrInvalidTransition
	default:
		return errs.ErrInvalidTransition
	}
	return nil
}
