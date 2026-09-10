package repository

import (
	"context"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"
	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
)

type AdminRepository interface {
	GetAdminByID(ctx context.Context, id int) (*entitie.AdminRes, error)
	GetAdminByUsername(ctx context.Context, username string) (*entitie.Admin, error)
	CreateAdmin(ctx context.Context, param db.CreateAdminParams) (int64, error)
	GetAllAdmins(ctx context.Context) ([]*entitie.AdminRes, error)
	UpdateAdmin(ctx context.Context, params db.UpdateAdminParams) error
	DeleteAdmin(ctx context.Context, params db.UpdateAdminIsActiveParams) error
}
