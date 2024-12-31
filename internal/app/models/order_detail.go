package models

import "github.com/google/uuid"

type OrderDetail struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    OrderID   uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
    ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
    Quantity  int32     `gorm:"type:int;not null" json:"quantity"`
    Price     int32     `gorm:"type:int;not null" json:"price"`
}

func (OrderDetail) TableName() string {
	return "order_detail"
}