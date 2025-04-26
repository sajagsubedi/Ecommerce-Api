package models

import "time"

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"not null" json:"name" validate:"required,min=2,max=100"`
	Email     string    `gorm:"unique;not null" json:"email" validate:"email,required"`
	Password  string    `gorm:"not null" json:"password" validate:"required,min=6"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
