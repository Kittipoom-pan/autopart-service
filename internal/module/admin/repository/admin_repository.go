package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
)

type AdminRepository interface {
	GetAdminByID(ctx context.Context, id int) (*entitie.AdminRes, error)
	GetAdminByUsername(ctx context.Context, username string) (*entitie.Admin, error)
	CreateAdmin(ctx context.Context, admin *entitie.AdminReq) (int64, error)
	GetAllAdmins(ctx context.Context) ([]*entitie.AdminRes, error)
	UpdateAdmin(ctx context.Context, id int, admin *entitie.AdminReq) error
	DeleteAdmin(ctx context.Context, id int) error
}
