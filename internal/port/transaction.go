package port

import (
	"context"

	"subscription-service/internal/app/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Transaction = (*pgxpool.Tx)(nil)

//go:generate mockgen -destination ../adapter/repo/mock/transaction_mock.go -package repo -source ./transaction.go

type TransactionController interface {
	BeginTx(ctx context.Context, isoLvl entity.IsolationLevel) (Transaction, error)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
