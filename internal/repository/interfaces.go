package repository

import (
	"context"
	"time"

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

type IOrderRepository interface {
	CreateFromCart(ctx context.Context, userID int) (*models.Order, error)
	GetByID(ctx context.Context, userID int, orderID uint) (*models.Order, error)
	ListByUser(ctx context.Context, userID int) ([]models.Order, error)
	ListAll(ctx context.Context) ([]models.Order, error)
	UpdateStatus(ctx context.Context, orderID uint, status string) error
}

type IPaymentRepository interface {
	FindByOrderID(ctx context.Context, userID int, orderID uint) (*models.Payment, error)
	FindByRazorpayOrderID(ctx context.Context, userID int, razorpayOrderID string) (*models.Payment, error)
	Create(ctx context.Context, payment *models.Payment) error
	MarkPaid(ctx context.Context, userID int, orderID uint, razorpayPaymentID string, signature string, paidAt time.Time) error
}
