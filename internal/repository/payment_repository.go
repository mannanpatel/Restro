package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mannanpatel/Restro/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrPaymentConflict = errors.New("payment conflicts with an existing payment")
var ErrOrderNotPayable = errors.New("order is not payable")

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &paymentRepository{db: db}
}

func (pr *paymentRepository) FindByOrderID(ctx context.Context, userID int, orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := pr.db.WithContext(ctx).Where("user_id = ? AND order_id = ?", userID, orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (pr *paymentRepository) FindByRazorpayOrderID(ctx context.Context, userID int, razorpayOrderID string) (*models.Payment, error) {
	var payment models.Payment
	err := pr.db.WithContext(ctx).Where("user_id = ? AND razorpay_order_id = ?", userID, razorpayOrderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (pr *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return pr.db.WithContext(ctx).Create(payment).Error
}

func (pr *paymentRepository) MarkPaid(ctx context.Context, userID int, orderID uint, razorpayPaymentID string, signature string, paidAt time.Time) error {
	return pr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payment models.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND order_id = ?", userID, orderID).
			First(&payment).Error; err != nil {
			return err
		}
		if payment.Status == "paid" {
			if payment.RazorpayPaymentID != nil && *payment.RazorpayPaymentID == razorpayPaymentID {
				return nil
			}
			return ErrPaymentConflict
		}
		if payment.Status != "created" {
			return ErrPaymentConflict
		}

		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&order).Error; err != nil {
			return err
		}
		if order.Status != "pending" {
			return ErrOrderNotPayable
		}

		if err := tx.Model(&payment).Updates(map[string]interface{}{
			"status":              "paid",
			"razorpay_payment_id": razorpayPaymentID,
			"razorpay_signature":  signature,
			"paid_at":             paidAt,
		}).Error; err != nil {
			return err
		}
		return tx.Model(&order).Updates(map[string]interface{}{
			"status":              "confirmed",
			"razorpay_payment_id": razorpayPaymentID,
			"razorpay_order_id":   payment.RazorpayOrderID,
		}).Error
	})
}
