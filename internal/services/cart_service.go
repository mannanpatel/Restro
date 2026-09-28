package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/repository"
)

type cartService struct {
	cartRepository    repository.ICartRepository
	productRepository repository.IProductRepository
}

func NewCartService(cartRepository repository.ICartRepository, productRepository repository.IProductRepository) ICartService {
	return &cartService{
		cartRepository:    cartRepository,
		productRepository: productRepository,
	}
}

func (cs *cartService) AddToCart(ctx context.Context, userID int, req models.AddCartItemRequest) error {
	// 1. Fetch the product from the database to ensure it exists
	product, err := cs.productRepository.GetProductById(ctx, req.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// 2. Validate Business Rules
	if !product.IsAvailable {
		return errors.New("this item is currently sold out or unavailable")
	}

	// 3. Pass the SERVER'S price, not the client's price, to the repository
	err = cs.cartRepository.AddOrUpdateCartItem(ctx, userID, req.ProductID, req.Quantity, product.Price)
	if err != nil {
		return err
	}
	return nil
}
