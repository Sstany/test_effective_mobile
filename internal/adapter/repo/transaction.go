package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"subscription-service/internal/app/entity"
	"subscription-service/internal/port"
)

var _ port.TransactionController = (*TransactionSQL)(nil)

type TransactionSQL struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewTransactionSQL(db *pgxpool.Pool, logger *zap.Logger) *TransactionSQL {
	return &TransactionSQL{
		db:     db,
		logger: logger,
	}
}

func (r *TransactionSQL) BeginTx(ctx context.Context, isoLvl entity.IsolationLevel) (port.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   toPgxIsoLevel(isoLvl),
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func toPgxIsoLevel(lvl entity.IsolationLevel) pgx.TxIsoLevel {
	switch lvl {
	case entity.RepeatableRead:
		return pgx.RepeatableRead
	case entity.Serializable:
		return pgx.Serializable
	default:
		return pgx.ReadCommitted
	}
}
