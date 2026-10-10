package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/repository"
	"gorm.io/gorm"
)

var ErrPaymentNotConfigured = errors.New("Razorpay credentials are not configured")
var ErrInvalidPaymentSignature = errors.New("invalid Razorpay signature")
var ErrProviderPaymentMismatch = errors.New("Razorpay payment does not match the order")

type razorpayGateway interface {
	KeyID() string
	CreateOrder(ctx context.Context, amount int64, receipt string) (*razorpayProviderOrder, error)
	FetchPayment(ctx context.Context, paymentID string) (*razorpayProviderPayment, error)
	CapturePayment(ctx context.Context, paymentID string, amount int64) (*razorpayProviderPayment, error)
}

type paymentService struct {
	payments repository.IPaymentRepository
	orders   IOrderService
	gateway  razorpayGateway
	secret   string
}

func NewPaymentService(payments repository.IPaymentRepository, orders IOrderService, gateway razorpayGateway, secret string) IPaymentService {
	return &paymentService{payments: payments, orders: orders, gateway: gateway, secret: secret}
}

func (s *paymentService) CreateOrder(ctx context.Context, userID int, orderID uint) (*models.CreatePaymentOrderResponse, error) {
	if userID <= 0 || orderID == 0 {
		return nil, errors.New("invalid payment request")
	}
	if s.secret == "" || s.gateway.KeyID() == "" {
		return nil, ErrPaymentNotConfigured
	}
	order, err := s.orders.GetByID(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}
	if order.Status != "pending" {
		return nil, repository.ErrOrderNotPayable
	}

	existing, err := s.payments.FindByOrderID(ctx, userID, orderID)
	if err == nil {
		if existing.Status != "created" {
			return nil, repository.ErrPaymentConflict
		}
		return paymentOrderResponse(s.gateway.KeyID(), existing)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	amount, err := amountInPaise(order.TotalPrice)
	if err != nil {
		return nil, err
	}
	providerOrder, err := s.gateway.CreateOrder(ctx, amount, fmt.Sprintf("restro_order_%d", order.ID))
	if err != nil {
		return nil, err
	}
	if providerOrder.Amount != amount || providerOrder.Currency != "INR" {
		return nil, errors.New("Razorpay order amount or currency did not match")
	}

	payment := &models.Payment{
		UserID:          userID,
		OrderID:         order.ID,
		Amount:          float64(amount) / 100,
		Status:          "created",
		RazorpayOrderID: providerOrder.ID,
	}
	if err := s.payments.Create(ctx, payment); err != nil {
		// A parallel checkout request may have persisted its provider order first.
		existing, lookupErr := s.payments.FindByOrderID(ctx, userID, orderID)
		if lookupErr == nil && existing.Status == "created" {
			return paymentOrderResponse(s.gateway.KeyID(), existing)
		}
		return nil, err
	}
	return paymentOrderResponse(s.gateway.KeyID(), payment)
}

func (s *paymentService) VerifyPayment(ctx context.Context, userID int, req models.VerifyPaymentRequest) error {
	if userID <= 0 || req.RazorpayOrderID == "" || req.RazorpayPaymentID == "" || req.RazorpaySignature == "" {
		return errors.New("invalid payment verification request")
	}
	if s.secret == "" {
		return ErrPaymentNotConfigured
	}
	if !validRazorpaySignature(s.secret, req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature) {
		return ErrInvalidPaymentSignature
	}

	payment, err := s.payments.FindByRazorpayOrderID(ctx, userID, req.RazorpayOrderID)
	if err != nil {
		return err
	}
	if payment.Status == "paid" {
		if payment.RazorpayPaymentID != nil && *payment.RazorpayPaymentID == req.RazorpayPaymentID {
			return nil
		}
		return repository.ErrPaymentConflict
	}
	if payment.Status != "created" {
		return repository.ErrPaymentConflict
	}

	providerPayment, err := s.gateway.FetchPayment(ctx, req.RazorpayPaymentID)
	if err != nil {
		return err
	}
	expectedAmount, err := amountInPaise(payment.Amount)
	if err != nil {
		return err
	}
	if providerPayment.ID != req.RazorpayPaymentID ||
		providerPayment.OrderID != payment.RazorpayOrderID ||
		providerPayment.Amount != expectedAmount ||
		providerPayment.Currency != "INR" {
		return ErrProviderPaymentMismatch
	}
	if providerPayment.Status == "authorized" {
		providerPayment, err = s.gateway.CapturePayment(ctx, req.RazorpayPaymentID, expectedAmount)
		if err != nil {
			return err
		}
	}
	if providerPayment.Status != "captured" {
		return errors.New("Razorpay payment is not captured")
	}

	return s.payments.MarkPaid(ctx, userID, payment.OrderID, req.RazorpayPaymentID, req.RazorpaySignature, time.Now().UTC())
}

func paymentOrderResponse(keyID string, payment *models.Payment) (*models.CreatePaymentOrderResponse, error) {
	amount, err := amountInPaise(payment.Amount)
	if err != nil {
		return nil, err
	}
	return &models.CreatePaymentOrderResponse{
		KeyID:           keyID,
		RazorpayOrderID: payment.RazorpayOrderID,
		Amount:          amount,
		Currency:        "INR",
		OrderID:         payment.OrderID,
	}, nil
}

func amountInPaise(amount float64) (int64, error) {
	if amount <= 0 || math.IsNaN(amount) || math.IsInf(amount, 0) {
		return 0, errors.New("order amount must be positive")
	}
	paise := math.Round(amount * 100)
	if math.Abs(amount*100-paise) > 0.000001 || paise > math.MaxInt64 {
		return 0, errors.New("order amount has invalid precision or is too large")
	}
	return int64(paise), nil
}

func validRazorpaySignature(secret, orderID, paymentID, signature string) bool {
	provided, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	fmt.Fprintf(mac, "%s|%s", orderID, paymentID)
	return hmac.Equal(provided, mac.Sum(nil))
}
