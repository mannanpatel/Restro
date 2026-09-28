package repository

import (
	"context"

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
		} else if err == gorm.ErrRecordNotFound {
			// Product doesn't exist → create new cart item
			cartItem = models.CartItem{
				CartID:    int(cart.ID),
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
