package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID     int        `json:"user_id" gorm:"not null;uniqueIndex"`
	TotalPrice float64    `json:"total_price" gorm:"type:decimal(10,2);default:0"`
	Items      []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model
	CartID    uint    `json:"cart_id" gorm:"not null;index"`
	ProductID int     `json:"product_id" gorm:"not null;index"` // Product ID
	Quantity  int     `json:"quantity" gorm:"not null;default:1"`
	Price     float64 `json:"price" gorm:"type:decimal(10,2);not null"` // snapshot price
}

type AddCartItemRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
