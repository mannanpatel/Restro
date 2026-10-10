package models

import (
	"time"

	"gorm.io/gorm"
)

type VerifyPaymentRequest struct {
	RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
	RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" binding:"required"`
}

type CreatePaymentOrderResponse struct {
	KeyID           string `json:"key_id"`
	RazorpayOrderID string `json:"razorpay_order_id"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	OrderID         uint   `json:"order_id"`
}

type Payment struct {
	gorm.Model
	UserID            int       `json:"user_id" gorm:"not null;index"`
	OrderID           uint      `json:"orderId" gorm:"uniqueIndex;not null"`
	Amount            float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status            string    `json:"status" gorm:"type:varchar(30);not null;default:'created'"`
	RazorpayPaymentID *string   `json:"razorpayPaymentId,omitempty" gorm:"type:varchar(191);uniqueIndex"`
	RazorpayOrderID   string    `json:"razorpayOrderId" gorm:"type:varchar(191);uniqueIndex;not null"`
	RazorpaySignature string    `json:"-" gorm:"type:varchar(255)"`
	PaidAt            time.Time `json:"paidAt,omitempty"`
}
