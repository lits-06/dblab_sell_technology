package models

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Price       uint      `gorm:"type:int;not null" json:"price"`
	Quantity    uint      `gorm:"type:int;default:0;not null" json:"quantity"`
	CategoryID  uint      `gorm:"type:int" json:"category_id"`
	Description *string   `gorm:"type:text" json:"description"`
	ImageURL    *string   `gorm:"type:varchar(255)" json:"image_url"`
}

func (Product) TableName() string {
	return "products"
}