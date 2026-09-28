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
	DeleteUser(id int) error
	ListUsers(ctx context.Context) ([]*models.User, error)
	UpdateUserRole(id int, Role string) error
}

type ICartRepository interface {
	AddOrUpdateCartItem(ctx context.Context, userID int, productID int, qty int, price float64) error
	GetCartByUserID(ctx context.Context, userID int) (*models.Cart, error)
	UpdateCartItemQuantity(ctx context.Context, userID int, itemID int, qty int) error
	DeleteCartItem(ctx context.Context, userID int, itemID int) error
	ClearCart(ctx context.Context, userID int) error
}
