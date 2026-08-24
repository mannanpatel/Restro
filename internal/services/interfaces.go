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
}

type IProductService interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductById(ctx context.Context, id int) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (int, error)
	UpdateProduct(ctx context.Context, id int, product *models.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id int) error
}
