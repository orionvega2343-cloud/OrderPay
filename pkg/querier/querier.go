package querier

import (
	"context"
	"database/sql"
)

// Querier - абстракция над *sqlx.DB и *sqlx.Tx,
// чтобы репозиторий мог одинаково работать
// как в рамках транзакции, так и без неё.
type Querier interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
