package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    Name      string    `gorm:"type:varchar(255)" json:"name"`
    TotalPrice int       `gorm:"type:int;not null" json:"total_price"`
    Status    string    `gorm:"type:varchar(50);default:'pending'" json:"status"`
    Phone     string    `gorm:"type:varchar(20);not null" json:"phone"`
    Address   string    `gorm:"type:text;not null" json:"address"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Order) TableName() string {
	return "orders"
}