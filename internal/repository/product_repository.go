package repository

import (
	"context"
	"fmt"

	"github.com/mannanpatel/Restro/internal/models"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}

// GetAllProducts fetches all products from database
func (pr *productRepository) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	var products []*models.Product

	result := pr.db.WithContext(ctx).Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", result.Error)
	}

	return products, nil
}

// GetProductById fetches a single product by ID
func (pr *productRepository) GetProductById(ctx context.Context, id int) (*models.Product, error) {
	var product *models.Product

	result := pr.db.WithContext(ctx).Where("id = ?", id).First(&product)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to fetch product: %w", result.Error)
	}

	return product, nil
}

// CreateProduct creates a new product
func (pr *productRepository) CreateProduct(ctx context.Context, product *models.Product) (int, error) {

	result := pr.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create product: %w", result.Error)
	}

	return int(product.ID), nil
}

// UpdateProduct updates an existing product
func (pr *productRepository) UpdateProduct(ctx context.Context, id int, product *models.UpdateProductRequest) error {
	updates := make(map[string]interface{})

	if product.Name != nil {
		updates["name"] = *product.Name
	}

	if product.Description != nil {
		updates["description"] = *product.Description
	}

	if product.Price != nil {
		updates["price"] = *product.Price
	}

	if product.Category != nil {
		updates["category"] = *product.Category
	}

	if product.ImageURL != nil {
		updates["image_url"] = *product.ImageURL
	}

	if product.IsAvailable != nil {
		updates["is_available"] = *product.IsAvailable
	}

	if product.Rating != nil {
		updates["rating"] = *product.Rating
	}

	if product.Type != nil {
		updates["type"] = *product.Type
	}

	if product.MealTime != nil {
		updates["meal_time"] = *product.MealTime
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result := pr.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update product: %w", result.Error)
	}

	return nil
}

// DeleteProduct deletes a product by ID
func (pr *productRepository) DeleteProduct(ctx context.Context, id int) error {
	result := pr.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Product{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
