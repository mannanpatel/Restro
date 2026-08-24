package repository

import (
	"context"

	"github.com/mannanpatel/Restro/internal/models"
)

// IProductRepository defines product repository interface
type IProductRepository interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductById(ctx context.Context, id int) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (int, error)
	UpdateProduct(ctx context.Context, id int, product *models.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id int) error
}

// IUserRepository defines user repository interface
type IUserRepository interface {
	RegisterUserWithProfile(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	GetUserById(id int) (*models.User, error)
	UpdateUser(id int, req *models.UpdateUserRequest) error
}
