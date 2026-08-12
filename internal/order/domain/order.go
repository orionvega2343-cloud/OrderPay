package domain

import (
	"errors"
	"fmt"
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
			return fmt.Errorf("invalid order status %s", next)
		}
	case "paid":
		switch next {
		case "completed", "cancelled":
			o.Status = next
		default:
			return fmt.Errorf("invalid order status %s", next)
		}
	case "cancelled", "completed":
		return errors.New("transition denied")
	default:
		return errors.New("Invalid status")
	}
	return nil
}
