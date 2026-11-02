package repository

import (
	"context"
	"database/sql"

	"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/payment/entity"
	"github.com/Kittipoom-pan/autopart-service/pkg/utils"
	"github.com/rs/zerolog"
)

type paymentRepository struct {
	conn    *sql.DB
	queries *db.Queries
	logger  zerolog.Logger
}

func NewPaymentRepository(conn *sql.DB, logger zerolog.Logger) PaymentRepository {
	return &paymentRepository{conn: conn, queries: db.New(conn), logger: logger}
}

func (r *paymentRepository) getQueries(ctx context.Context) *db.Queries {
	if tx, ok := database.GetTx(ctx); ok {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *paymentRepository) CreatePayment(ctx context.Context, in *entity.CreatePaymentInput) error {
	_, err := r.getQueries(ctx).CreatePayment(ctx, db.CreatePaymentParams{
		OrderID:                 in.OrderID,
		Amount:                  in.Amount,
		CustomerPaymentMethodID: utils.Int32PtrToNullInt32(in.CustomerPaymentMethodID),
		PaymentMethodTypeID:     in.PaymentMethodTypeID,
		Status:                  db.PaymentStatusPending,
		CreatedBy:               utils.Int32PtrToNullInt32(in.CreatedBy),
	})
	if err != nil {
		r.logger.Error().Err(err).Int64("order_id", in.OrderID).Msg("failed to create payment")
		return dberr.HandleMySQLError(err)
	}
	return nil
}
