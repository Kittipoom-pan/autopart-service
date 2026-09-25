package helper

import (
	"errors"

	"github.com/Kittipoom-pan/autopart-service/internal/common"
	dberr "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/dberror"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

// converts database errors to API errors
func MapDBErrorToAPIError(err error, resource string) error {
	if err == nil {
		return nil
	}

	var dbErr *dberr.DBError
	if errors.As(err, &dbErr) {
		if dberr.IsDuplicateKey(dbErr) {
			return customerror.NewAPIError(common.StatusConflict, resource+" already exists (duplicate key)")
		}
		if dberr.IsForeignKeyError(dbErr) {
			return customerror.NewAPIError(common.StatusBadRequest, "Invalid reference: "+dbErr.Detail)
		}
		if dberr.IsRecordNotFound(dbErr) {
			return customerror.NewNotFoundError(resource)
		}
		if errors.Is(dbErr.Err, dberr.ErrDatabaseInternal) {
			return customerror.NewAPIError(common.StatusError, "database operation failed")
		}
	}

	return err
}
