package services

import (
	"context"

	"github.com/mannanpatel/Restro/internal/models"
)

type IUserService interface {
	Register(user *models.User) error
	Login(req models.UserLoginRequest) (*models.UserLoginResponse, error)
	GetUserById(id int) (*models.User, error)
	UpdateUser(id int, req *models.UpdateUserRequest) error
	DeleteUser(id int) error
	ListUsers(ctx context.Context) ([]*models.User, error)
	UpdateUserRole(id int, Role string) error
}

type IProductService interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductById(ctx context.Context, id int) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (int, error)
	UpdateProduct(ctx context.Context, id int, product *models.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id int) error
}

type ICartService interface {
	AddToCart(ctx context.Context, userID int, req models.AddCartItemRequest) error
	GetCart(ctx context.Context, userID int) (*models.Cart, error)
	UpdateCartItemQuantity(ctx context.Context, userID int, itemID int, qty int) error
	RemoveCartItem(ctx context.Context, userID int, itemID int) error
	ClearCart(ctx context.Context, userID int) error
}
