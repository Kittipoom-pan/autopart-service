package entitie

type Customer struct {
	ID       uint32 `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
