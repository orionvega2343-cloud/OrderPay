package repository

import (
	"OrderPay/internal/payment/domain"
	"OrderPay/pkg/transaction"
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type PaymentRepo interface {
	CreatePayment(ctx context.Context, m *domain.Payment) (*domain.Payment, error)
	GetPaymentById(ctx context.Context, id int) (*domain.Payment, error)
	GetAllPayments(ctx context.Context) ([]domain.Payment, error)
	UpdatePayment(ctx context.Context, m *domain.Payment) error
	DeletePayment(ctx context.Context, id int) error
	GetPaymentTotalAmount(ctx context.Context, userId, status string, from, to time.Time) (int, error)
}

type Querier interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
type PaymentRepoImpl struct {
	db *sqlx.DB
}

func NewPaymentRepo(db *sqlx.DB) *PaymentRepoImpl {
	return &PaymentRepoImpl{db: db}
}

func (r *PaymentRepoImpl) CreatePayment(ctx context.Context, m *domain.Payment) (*domain.Payment, error) {
	//Получаем значение из контекста
	//и проверяем с помощью assertion type
	var q Querier
	tx, ok := transaction.ExtractTx(ctx)
	query := `INSERT INTO payments(order_id, amount, status, method) VALUES($1, $2, $3, $4) RETURNING id`
	if ok {
		q = tx
	} else {
		q = r.db
	}
	err := q.GetContext(ctx, &m.Id, query, m.OrderId, m.Amount, m.Status, m.Method)
	if err != nil {
		slog.Error("Failed to create payment", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *PaymentRepoImpl) GetPaymentById(ctx context.Context, id int) (*domain.Payment, error) {
	var m domain.Payment
	query := `SELECT id, order_id, amount, status, method FROM payments WHERE id = $1`
	err := r.db.GetContext(ctx, &m, query, id)
	if err != nil {
		slog.Error("Failed to get payment by id", "error", err)
		return nil, err
	}
	return &m, nil
}

func (r *PaymentRepoImpl) GetAllPayments(ctx context.Context) ([]domain.Payment, error) {
	var m []domain.Payment
	query := `SELECT id, order_id, amount, status, method FROM payments`
	err := r.db.SelectContext(ctx, &m, query)
	if err != nil {
		slog.Error("Failed to get all payments", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *PaymentRepoImpl) UpdatePayment(ctx context.Context, m *domain.Payment) error {
	tx, ok := transaction.ExtractTx(ctx)
	var q Querier
	query := `UPDATE payments SET status = $1, method = $2 WHERE id = $3`
	if ok {
		q = tx
	} else {
		q = r.db
	}
	_, err := q.ExecContext(ctx, query, m.Status, m.Method, m.Id)
	if err != nil {
		slog.Error("Failed to update payment", "error", err)
		return err
	}
	return nil
}

func (r *PaymentRepoImpl) DeletePayment(ctx context.Context, id int) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete payment", "error", err)
		return err
	}
	return nil
}

func (r *PaymentRepoImpl) GetPaymentTotalAmount(ctx context.Context, userId, status string, from, to time.Time) (int, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM payments JOIN orders ON(payments.order_id = orders.id) WHERE orders.user_id = $1 AND payments.status = $2 AND payments.created_at BETWEEN $3 AND $4`
	var total int
	err := r.db.GetContext(ctx, &total, query, userId, status, from, to)
	if err != nil {
		slog.Error("Failed to get total amount of payments", "error", err)
		return 0, err
	}
	return total, nil
}
