package entitie

import "time"

type Customer struct {
	ID          uint32    `json:"id" db:"id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password" db:"password"`
	BirthDate   time.Time `json:"birth_date" db:"birth_date"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
}
