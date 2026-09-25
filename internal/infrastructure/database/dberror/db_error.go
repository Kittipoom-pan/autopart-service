package dberror

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const (
	ErrDuplicateEntry   uint16 = 1062
	ErrForeignKeyFailed uint16 = 1452
)

var (
	ErrDuplicateKey     = errors.New("duplicate key error")
	ErrForeignKey       = errors.New("foreign key constraint failed")
	ErrRecordNotFound   = errors.New("record not found")
	ErrDatabaseInternal = errors.New("database internal error")
)

type DBError struct {
	Err    error
	Detail string
}

func (e *DBError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Err.Error(), e.Detail)
	}
	return e.Err.Error()
}

func (e *DBError) Unwrap() error {
	return e.Err
}

// NewDBError creates a new DBError
func NewDBError(err error, detail string) *DBError {
	return &DBError{
		Err:    err,
		Detail: detail,
	}
}

// HandleMySQLError converts MySQL errors to domain errors
func HandleMySQLError(err error) error {
	if err == nil {
		return nil
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case ErrDuplicateEntry:
			return NewDBError(ErrDuplicateKey, mysqlErr.Message)
		case ErrForeignKeyFailed:
			return NewDBError(ErrForeignKey, mysqlErr.Message)
		}
	}

	return NewDBError(ErrDatabaseInternal, err.Error())
}

// IsDuplicateKey checks if error is a duplicate key error
func IsDuplicateKey(err error) bool {
	var dbErr *DBError
	if errors.As(err, &dbErr) {
		return errors.Is(dbErr.Err, ErrDuplicateKey)
	}
	return false
}

// IsForeignKeyError checks if error is a foreign key error
func IsForeignKeyError(err error) bool {
	var dbErr *DBError
	if errors.As(err, &dbErr) {
		return errors.Is(dbErr.Err, ErrForeignKey)
	}
	return false
}

// IsRecordNotFound checks if error is a record not found error
func IsRecordNotFound(err error) bool {
	var dbErr *DBError
	if errors.As(err, &dbErr) {
		return errors.Is(dbErr.Err, ErrRecordNotFound)
	}
	return false
}
