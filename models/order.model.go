package models

import (
	"time"
)

const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
)

type Order struct {
	ID              uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          uint            `json:"user_id" gorm:"not null"`
	User            User            `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Items           []OrderItem     `json:"items" gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
	TotalAmount     float64         `json:"total_amount" gorm:"not null" validate:"required,gt=0"`
	Status          string          `json:"status" gorm:"type:varchar(20);default:'pending'" validate:"oneof=pending processing shipped delivered cancelled"`
	ShippingAddress ShippingAddress `json:"shipping_address" gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ContactNumber   string          `json:"contact_number" gorm:"type:varchar(10);not null" validate:"required"`
	CreatedAt       time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"Product" gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity  int     `json:"quantity" gorm:"not null" validate:"required,min=1"`
	Price     float64 `json:"price" gorm:"not null" validate:"required,gte=0"`
}
