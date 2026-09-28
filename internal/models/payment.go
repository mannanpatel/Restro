package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentRequestBody struct {
	Amount int  `json:"amount" binding:"required,gt=0"`
	UserID uint `json:"userId" binding:"required"`
}

type VerifyPaymentBody struct {
	RazorpayPaymentID string     `json:"razorpay_payment_id" binding:"required"`
	RazorpayOrderID   string     `json:"razorpay_order_id" binding:"required"`
	RazorpaySignature string     `json:"razorpay_signature" binding:"required"`
	CartItems         []CartItem `json:"cartItems"`
	TotalPrice        int        `json:"totalPrice" binding:"required,gt=0"`
	UserID            uint       `json:"userId" binding:"required"`
	DateTime          time.Time  `json:"dateTime"`
}

type PaymentResponseBody struct {
	PaymentID string `json:"paymentId"`
	Message   string `json:"message"`
}

type Payment struct {
	gorm.Model
	UserID            uint      `json:"userId" gorm:"not null;index"`
	OrderID           uint      `json:"orderId" gorm:"index"`
	Amount            float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status            string    `json:"status" gorm:"type:varchar(30);not null;default:'created'"`
	RazorpayPaymentID string    `json:"razorpayPaymentId" gorm:"type:varchar(191);uniqueIndex"`
	RazorpayOrderID   string    `json:"razorpayOrderId" gorm:"type:varchar(191);uniqueIndex;not null"`
	RazorpaySignature string    `json:"razorpaySignature" gorm:"type:varchar(255)"`
	PaidAt            time.Time `json:"paidAt"`
}
