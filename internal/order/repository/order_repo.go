package repository

import (
	"OrderPay/internal/order/domain"
	"OrderPay/pkg/transaction"
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, m *domain.Order, items []domain.OrderItem) (*domain.Order, error)
	GetOrderByID(ctx context.Context, id int) (domain.Order, error)
	GetAllOrders(ctx context.Context) ([]domain.Order, error)
	UpdateOrder(ctx context.Context, m *domain.Order) error
	DeleteOrder(ctx context.Context, id int) error
}

// Querier - отдельный интерфейс,
// для устранения дублирования кода внутри CreateOrder
type Querier interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type OrderRepoImpl struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepoImpl {
	return &OrderRepoImpl{db: db}
}

//определение транзакции, для метода через ExtractTx
//если транзакция нужна, репозиторий делает запрос через tx,
//иначе через r.db

func (r *OrderRepoImpl) CreateOrder(ctx context.Context, m *domain.Order, items []domain.OrderItem) (*domain.Order, error) {
	var q Querier
	tx, ok := transaction.ExtractTx(ctx)
	query := `INSERT INTO orders(user_id, status, total_amount) VALUES($1, $2, $3) RETURNING id`
	queryItems := `INSERT INTO order_items(order_id, product_name, quantity, price_per_unit) VALUES($1, $2, $3, $4) RETURNING id`
	if ok {
		q = tx
	} else {
		q = r.db
	}
	err := q.GetContext(ctx, &m.Id, query, m.UserId, m.Status, m.TotalAmount)
	if err != nil {
		slog.Error("failed inserting order", "error", err)
		return nil, err
	}
	for _, item := range items {
		err = q.GetContext(ctx, &item.Id, queryItems, m.Id, item.ProductName, item.Quantity, item.PricePerUnit)
		if err != nil {
			return nil, err
		}
	}
	return m, nil

}

func (r *OrderRepoImpl) GetOrderByID(ctx context.Context, id int) (domain.Order, error) {
	var m domain.Order
	query := `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE id = $1`
	err := r.db.GetContext(ctx, &m, query, id)
	if err != nil {
		slog.Error("failed getting order", "error", err)
		return m, err
	}
	return m, nil
}

func (r *OrderRepoImpl) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	var m []domain.Order
	query := `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders`
	err := r.db.SelectContext(ctx, &m, query)
	if err != nil {
		slog.Error("failed getting all orders", "error", err)
		return m, err
	}
	return m, nil
}

func (r *OrderRepoImpl) UpdateOrder(ctx context.Context, m *domain.Order) error {
	var q Querier()
	tx, ok := transaction.ExtractTx(ctx)
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	if ok {
		q = tx
	} else {
		q = r.db
	}
	_, err := q.ExecContext(ctx, query, m.Status, m.Id)
	if err != nil {
		slog.Error("failed updating order", "error", err)
		return err
	}
	return nil
}

func (r *OrderRepoImpl) DeleteOrder(ctx context.Context, id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("failed deleting order", "error", err)
		return err
	}
	return nil
}
