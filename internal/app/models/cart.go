package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID	`gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID	`gorm:"type:uuid;not null" json:"user_id"`
	ProductID uuid.UUID	`gorm:"type:uuid;not null" json:"product_id"`
	Quantity  int       `gorm:"type:int;not null" json:"quantity"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at" `
}

func (Cart) TableName() string {
	return "cart"
}