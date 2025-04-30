package models

import "time"

type ShippingAddress struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   uint      `json:"order_id" gorm:"not null"`
	Street    string    `json:"street" gorm:"type:varchar(255);not null" validate:"required"`
	City      string    `json:"city" gorm:"type:varchar(100);not null" validate:"required"`
	State     string    `json:"state" gorm:"type:varchar(100);not null" validate:"required"`
	ZipCode   string    `json:"zip_code" gorm:"type:varchar(20);not null" validate:"required"`
	Country   string    `json:"country" gorm:"type:varchar(100);not null" validate:"required"`
	Notes     string    `json:"notes" gorm:"type:text;optional"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
