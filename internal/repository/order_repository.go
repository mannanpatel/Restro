package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mannanpatel/Restro/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrCartEmpty = errors.New("cart is empty")
var ErrProductUnavailable = errors.New("a cart product is unavailable")

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &orderRepository{db: db}
}

func (or *orderRepository) CreateFromCart(ctx context.Context, userID int) (*models.Order, error) {
	var order models.Order
	err := or.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&cart).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCartEmpty
			}
			return err
		}

		var cartItems []models.CartItem
		if err := tx.Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
			return err
		}
		if len(cartItems) == 0 {
			return ErrCartEmpty
		}

		order = models.Order{
			UserID: userID,
			Status: "pending",
			Items:  make([]models.OrderItem, 0, len(cartItems)),
		}
		for _, cartItem := range cartItems {
			var product models.Product
			if err := tx.Where("id = ? AND is_available = ?", cartItem.ProductID, true).First(&product).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("%w: product %d", ErrProductUnavailable, cartItem.ProductID)
				}
				return err
			}
			order.TotalPrice += product.Price * float64(cartItem.Quantity)
			order.Items = append(order.Items, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				UnitPrice: product.Price,
			})
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}
		return tx.Model(&cart).Update("total_price", 0).Error
	})
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *orderRepository) GetByID(ctx context.Context, userID int, orderID uint) (*models.Order, error) {
	var order models.Order
	err := or.db.WithContext(ctx).Preload("Items").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *orderRepository) ListByUser(ctx context.Context, userID int) ([]models.Order, error) {
	var orders []models.Order
	err := or.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (or *orderRepository) ListAll(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	err := or.db.WithContext(ctx).Preload("Items").Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (or *orderRepository) UpdateStatus(ctx context.Context, orderID uint, status string) error {
	var order models.Order
	if err := or.db.WithContext(ctx).Select("id").First(&order, orderID).Error; err != nil {
		return err
	}
	return or.db.WithContext(ctx).Model(&order).Update("status", status).Error
}
