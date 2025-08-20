package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
)

type AdminUsecase interface {
	GetAdminByID(ctx context.Context, id int) (*entitie.AdminRes, error)
	GetAllAdmins(ctx context.Context) ([]*entitie.AdminRes, error)
	CreateAdmin(ctx context.Context, user *entitie.AdminReq) (int64, error)
	UpdateAdmin(ctx context.Context, id int, user *entitie.AdminReq) error
	DeleteAdmin(ctx context.Context, id int) error
}
