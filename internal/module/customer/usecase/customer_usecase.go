package usecase

import (
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/entitie"
	"github.com/Kittipoom-pan/autopart-service/internal/module/customer/repository"
)

type CustomerUsecase interface {
	GetCustomerByID(id int) (*entitie.Customer, error)
	CreateCustomer(user *entitie.Customer) (int64, error)
}

type CustomerUsecaseImpl struct {
	repo repository.CustomerRepository
}

func NewCustomerUsecase(repo repository.CustomerRepository) CustomerUsecase {
	return &CustomerUsecaseImpl{repo: repo}
}

func (u *CustomerUsecaseImpl) GetCustomerByID(id int) (*entitie.Customer, error) {
	customer, err := u.repo.GetCustomerByID(id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (u *CustomerUsecaseImpl) CreateCustomer(user *entitie.Customer) (int64, error) {
	// params := db.CreateCustomerParams{
	// 	FirstName:   sql.NullString{String: user.FirstName, Valid: user.FirstName != ""},
	// 	LastName:    sql.NullString{String: user.LastName, Valid: user.LastName != ""},
	// 	Username:    user.Username,
	// 	Email:       user.Email,
	// 	Password:    sql.NullString{String: user.Password, Valid: user.Password != ""},
	// 	BirthDate:   sql.NullTime{Time: user.BirthDate, Valid: !user.BirthDate.IsZero()},
	// 	PhoneNumber: sql.NullString{String: user.PhoneNumber, Valid: user.PhoneNumber != ""},
	// 	CreatedAt:   user.CreatedAt,
	// 	CreatedBy:   sql.NullString{String: user.CreatedBy, Valid: user.CreatedBy != ""},
	// 	UpdatedAt:   sql.NullTime{Time: user.UpdatedAt, Valid: !user.UpdatedAt.IsZero()},
	// 	UpdatedBy:   sql.NullString{String: user.UpdatedBy, Valid: user.UpdatedBy != ""},
	// }
	customerID, err := u.repo.CreateCustomer(user)
	if err != nil {
		return 0, err
	}
	return customerID, nil
}
