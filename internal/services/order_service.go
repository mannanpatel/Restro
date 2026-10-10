package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/repository"
)

type orderService struct {
	repository repository.IOrderRepository
}

func NewOrderService(repo repository.IOrderRepository) IOrderService {
	return &orderService{repository: repo}
}

func (s *orderService) CreateFromCart(ctx context.Context, userID int) (*models.Order, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repository.CreateFromCart(ctx, userID)
}

func (s *orderService) GetByID(ctx context.Context, userID int, orderID uint) (*models.Order, error) {
	if userID <= 0 || orderID == 0 {
		return nil, errors.New("invalid order request")
	}
	return s.repository.GetByID(ctx, userID, orderID)
}

func (s *orderService) ListByUser(ctx context.Context, userID int) ([]models.Order, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repository.ListByUser(ctx, userID)
}

func (s *orderService) ListAll(ctx context.Context) ([]models.Order, error) {
	return s.repository.ListAll(ctx)
}

func (s *orderService) UpdateStatus(ctx context.Context, orderID uint, status string) error {
	if orderID == 0 {
		return errors.New("invalid order id")
	}
	switch status {
	case "pending", "confirmed", "preparing", "ready", "completed", "cancelled":
	default:
		return fmt.Errorf("invalid order status")
	}
	return s.repository.UpdateStatus(ctx, orderID, status)
}
