package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entity"
)

type AuthUsecase interface {
	Login(ctx context.Context, request *entity.LoginRequest) (*entity.LoginResponse, error)
}
