package transaction

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type ctxKey string

const txKey ctxKey = "tx"

type Transactor struct {
	db *sqlx.DB
}

func NewTransactor(db *sqlx.DB) *Transactor {
	return &Transactor{db: db}
}

//UoW (Unit of Work) - паттерн, который гарантирует атомарность операций,
//затрагивающих несколько Repository в рамках одной бизнес-транзакции:
//либо все SQL-изменения внутри блока применяются успешно (commit),
//либо откатываются целиком при любой ошибке (rollback). Repository,
//вызванные внутри одного UoW-блока, работают с одной и той же SQL-транзакцией через общий контекст,
//вместо того чтобы каждый открывал и коммитил свою отдельную транзакцию независимо.

func (t *Transactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	trx, err := t.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		//Recover не дает транзакции зависнуть в бд
		if p := recover(); p != nil {
			_ = trx.Rollback()
			slog.Error("transaction panicked", "error", p)
			panic(p)
		}

		if err != nil {
			_ = trx.Rollback()

		} else {
			_ = trx.Commit()
		}

	}()

	ctx = context.WithValue(ctx, txKey, trx)

	err = fn(ctx)
	return err
}

func (t *Transactor) ExtractTx(ctx context.Context) (*sqlx.Tx, bool) {
	if v, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return v, true
	}
	return nil, false
}
