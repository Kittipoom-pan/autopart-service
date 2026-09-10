package repository

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
)

type AdminRepository interface {
	GetAdminByID(ctx context.Context, id int) (*entity.AdminRes, error)
	GetAdminByUsername(ctx context.Context, username string) (*entity.Admin, error)
	CreateAdmin(ctx context.Context, admin *entity.AdminReq, createdBy *int) (int64, error)
	GetAllAdmins(ctx context.Context) ([]*entity.AdminRes, error)
	UpdateAdmin(ctx context.Context, adminID int, admin *entity.AdminReq, updatedBy int) error
	DeleteAdmin(ctx context.Context, adminID int, updatedBy int) error
}
