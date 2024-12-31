package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email    string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password string    `gorm:"type:varchar(255);not null" json:"password"`
	Name     string    `gorm:"type:varchar(255)" json:"name"`
	Phone    *string   `gorm:"type:varchar(20)" json:"phone"`
	Address  *string   `gorm:"type:text" json:"address"`
	Role     string    `gorm:"type:varchar(50);not null" json:"role"`
	Avatar	 *string   `gorm:"type:varchar(255)" json:"avatar"`
}

func (User) TableName() string {
	return "users"
}