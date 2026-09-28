package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID            int         `json:"user_id" gorm:"not null;index"`
	Items             []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	TotalPrice        float64     `json:"totalPrice" gorm:"type:decimal(10,2);not null;default:0"`
	Status            string      `json:"status" gorm:"type:varchar(30);not null;default:'pending'"`
	RazorpayPaymentID string      `json:"razorpayPaymentId" gorm:"type:varchar(191);index"`
	RazorpayOrderID   string      `json:"razorpayOrderId" gorm:"type:varchar(191);index"`
}

type OrderItem struct {
	gorm.Model
	OrderID   int     `json:"order_id" gorm:"not null;index"`
	ProductID int     `json:"product_id" gorm:"not null;index"` // Product ID
	Quantity  int     `json:"quantity" gorm:"not null;default:1"`
	UnitPrice float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"` // snapshot price
}
