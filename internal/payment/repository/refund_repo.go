package repository

import (
	"OrderPay/internal/payment/domain"
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type RefundRepo interface {
	CreateRefund(ctx context.Context, m *domain.Refund) (*domain.Refund, error)
	GetAllRefund(ctx context.Context) ([]domain.Refund, error)
}

type RefundRepoImpl struct {
	db *sqlx.DB
}

func NewRefundRepo(db *sqlx.DB) *RefundRepoImpl {
	return &RefundRepoImpl{db: db}
}

func (r *RefundRepoImpl) CreateRefund(ctx context.Context, m *domain.Refund) (*domain.Refund, error) {
	query := `INSERT INTO refunds(payment_id, amount, status,  reason) VALUES($1, $2, $3, $4) RETURNING id`
	err := r.db.GetContext(ctx, &m.Id, query, m.PaymentId, m.Amount, m.Status, m.Reason)
	if err != nil {
		slog.Error("Failed to create refund", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *RefundRepoImpl) GetAllRefund(ctx context.Context) ([]domain.Refund, error) {
	var refunds []domain.Refund
	query := `SELECT id, payment_id, amount, reason, status, created_at FROM refunds`
	err := r.db.SelectContext(ctx, &refunds, query)
	if err != nil {
		slog.Error("Failed to get all refunds", "error", err)
		return nil, err
	}
	return refunds, nil
}
