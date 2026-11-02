package database

import (
	"context"
	"database/sql"

	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	"github.com/rs/zerolog"
)

type txKey struct{}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type sqlTxManager struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewTxManager(db *sql.DB, logger zerolog.Logger) TxManager {
	return &sqlTxManager{
		db:     db,
		logger: logger,
	}
}

func (m *sqlTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := GetTx(ctx); ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return dberr.NewDBError(dberr.ErrDatabaseInternal, "begin tx: "+err.Error())
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			m.logger.Error().Err(rbErr).Msg("rollback failed")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return dberr.NewDBError(dberr.ErrDatabaseInternal, "commit tx: "+err.Error())
	}
	return nil
}

// Helper ให้ Repository ดึง Tx หรือ DB connection ออกมาใช้
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}
