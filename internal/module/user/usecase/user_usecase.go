package usecase

import (
	"github.com/Kittipoom-pan/autopart-service/internal/module/user/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/user/repository"
)

type UsersUsecase interface {
	GetUser(id int) (*entitie.User, error)
	CreateUser(user *entitie.User) error
}

type UsersUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUsersUsecase(repo repository.UserRepository) UsersUsecase {
	return &UsersUsecaseImpl{repo: repo}
}

func (u *UsersUsecaseImpl) GetUser(id int) (*entitie.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *UsersUsecaseImpl) CreateUser(user *entitie.User) error {
	return u.repo.CreateUser(user)
}
