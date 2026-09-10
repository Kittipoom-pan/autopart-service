package entity

type Customer struct {
	ID       uint32 `json:"id" db:"id"`
	Uuid     string `json:"uuid"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
