package usecase

import (
	"context"

	"github.com/Kittipoom-pan/autopart-service/internal/module/admin/entitie"
)

type AuthUsecase interface {
	Login(ctx context.Context, request *entitie.LoginRequest) (*entitie.LoginResponse, error)
}
