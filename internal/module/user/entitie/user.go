package entitie

import "time"

type (
	User struct {
		Id        uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
		Amount    uint32    `json:"amount"`
		CreatedAt time.Time `json:"createdAt"`
	}
)
