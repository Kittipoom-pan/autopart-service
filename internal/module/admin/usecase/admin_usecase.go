package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
)

type AdminUsecase interface {
	GetAdminByID(ctx context.Context, id int) (*entity.AdminRes, error)
	GetAllAdmins(ctx context.Context) ([]*entity.AdminRes, error)
	CreateAdmin(ctx context.Context, user *entity.AdminReq, userID *int) (int64, error)
	UpdateAdmin(ctx context.Context, adminID int, userID int, user *entity.AdminReq) error
	DeleteAdmin(ctx context.Context, adminID int, userID int) error
}
