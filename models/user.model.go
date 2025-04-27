package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"not null" json:"name" validate:"required,min=2,max=100"`
	Email     string    `gorm:"unique;not null" json:"email" validate:"email,required"`
	Password  string    `gorm:"not null" json:"password" validate:"required,min=6"`
	Role      string    `gorm:"type:varchar(20);default:user" json:"role,omitempty" validate:"omitempty,oneof=user admin"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (user *User) HashPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (user *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
