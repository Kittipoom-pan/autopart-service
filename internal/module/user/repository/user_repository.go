package repository

import (
	"database/sql"
	//"github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database"
	"github.com/Kittipoom-pan/autopart-service/internal/module/user/entitie"
)

type UserRepository interface {
	GetUserByID(id int) (*entitie.User, error)
	CreateUser(user *entitie.User) error
	// ... other methods ...
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetUserByID(id int) (*entitie.User, error) {
	// ... implementation ...
	return nil, nil
}

func (r *UserRepositoryImpl) CreateUser(user *entitie.User) error {
	// ... implementation ...
	return nil
}
