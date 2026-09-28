package repository

import (
	"context"
	"errors"

	"github.com/mannanpatel/Restro/internal/models"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) ICartRepository {
	return &cartRepository{
		db: db,
	}
}

// func (cr *cartRepository) AddItem(userID uint, itemID uint, qty int, price float64) error {
// 	var cart models.Cart
// }

func (cr *cartRepository) AddOrUpdateCartItem(ctx context.Context, userID int, productID int, qty int, price float64) error {
	return cr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Find user's cart or create a new one
		var cart models.Cart
		err := tx.Where("user_id = ?", userID).FirstOrCreate(&cart, models.Cart{UserID: userID}).Error
		if err != nil {
			return err
		}

		// 2. Check if product already exists in cart
		var cartItem models.CartItem
		err = tx.Where("cart_id =? AND product_id =?", cart.ID, productID).First(&cartItem).Error
		if err == nil {
			// Product already exists → increase quantity
			cartItem.Quantity += qty
			err := tx.Save(&cartItem).Error
			if err != nil {
				return err
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Product doesn't exist → create new cart item
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: productID,
				Quantity:  qty,
				Price:     price,
			}
			err := tx.Create(&cartItem).Error
			if err != nil {
				return err
			}
		} else {
			return err
		}

		// 3. Recalculate cart total
		var total float64
		err = tx.Model(&models.CartItem{}).Where("cart_id = ?", cart.ID).Select("COALESCE(SUM(quantity * price), 0)").Scan(&total).Error
		if err != nil {
			return err
		}

		// 4. Update cart total
		return tx.Model(&cart).Update("total_price", total).Error
	})
}

func (cr *cartRepository) GetCartByUserID(ctx context.Context, userID int) (*models.Cart, error) {
	var cart models.Cart
	err := cr.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &models.Cart{UserID: userID, Items: []models.CartItem{}}, nil
	}
	if err != nil {
		return nil, err
	}
	if cart.Items == nil {
		cart.Items = []models.CartItem{}
	}
	return &cart, nil
}

func (cr *cartRepository) UpdateCartItemQuantity(ctx context.Context, userID int, itemID int, qty int) error {
	return cr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cart, err := findCartForUser(tx, userID)
		if err != nil {
			return err
		}
		var item models.CartItem
		if err := tx.Where("id = ? AND cart_id = ?", itemID, cart.ID).First(&item).Error; err != nil {
			return err
		}
		if err := tx.Model(&item).Update("quantity", qty).Error; err != nil {
			return err
		}
		return recalculateCartTotal(tx, cart.ID)
	})
}

func (cr *cartRepository) DeleteCartItem(ctx context.Context, userID int, itemID int) error {
	return cr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cart, err := findCartForUser(tx, userID)
		if err != nil {
			return err
		}
		result := tx.Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&models.CartItem{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return recalculateCartTotal(tx, cart.ID)
	})
}

func (cr *cartRepository) ClearCart(ctx context.Context, userID int) error {
	return cr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cart, err := findCartForUser(tx, userID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}
		return tx.Model(cart).Update("total_price", 0).Error
	})
}

func findCartForUser(tx *gorm.DB, userID int) (*models.Cart, error) {
	var cart models.Cart
	if err := tx.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func recalculateCartTotal(tx *gorm.DB, cartID uint) error {
	var total float64
	if err := tx.Model(&models.CartItem{}).
		Where("cart_id = ?", cartID).
		Select("COALESCE(SUM(quantity * price), 0)").
		Scan(&total).Error; err != nil {
		return err
	}
	return tx.Model(&models.Cart{}).Where("id = ?", cartID).Update("total_price", total).Error
}
