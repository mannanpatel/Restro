package services

import (
	"context"
	"fmt"

	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/repository"
)

type productService struct {
	repository repository.IProductRepository
}

func NewProductService(repo repository.IProductRepository) IProductService {
	return &productService{repository: repo}
}

// GetAllProducts fetches all products
func (ps *productService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	products, err := ps.repository.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products, nil
}

// GetProductById fetches a single product by ID
func (ps *productService) GetProductById(ctx context.Context, id int) (*models.Product, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid product id")
	}

	product, err := ps.repository.GetProductById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	return product, nil
}

// CreateProduct creates a new product
func (ps *productService) CreateProduct(ctx context.Context, product *models.Product) (int, error) {
	if product == nil {
		return 0, fmt.Errorf("product cannot be nil")
	}

	productID, err := ps.repository.CreateProduct(ctx, product)
	if err != nil {
		return 0, fmt.Errorf("failed to create product: %w", err)
	}
	return productID, nil
}

// UpdateProduct updates an existing product
func (ps *productService) UpdateProduct(ctx context.Context, id int, product *models.UpdateProductRequest) error {
	if id == 0 {
		return fmt.Errorf("invalid product id")
	}
	if product == nil {
		return fmt.Errorf("product cannot be nil")
	}

	err := ps.repository.UpdateProduct(ctx, id, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// DeleteProduct deletes a product by ID
func (ps *productService) DeleteProduct(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("invalid product id")
	}

	err := ps.repository.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
