package routes

import (
	"github.com/mannanpatel/Restro/internal/controller"
	"github.com/mannanpatel/Restro/internal/repository"
	"github.com/mannanpatel/Restro/internal/services"
	"gorm.io/gorm"
)

// CONTROLLERS
var (
	authController    *controller.AuthController
	productController *controller.ProductController
	userController    *controller.UserController
	cartController    *controller.CartController
)

// SERVICES
var (
	userService    services.IUserService
	productService services.IProductService
	cartService    services.ICartService
)

// REPOSITORIES
var (
	userRepository    repository.IUserRepository
	productRepository repository.IProductRepository
	cartRepository    repository.ICartRepository
)

// initDependencies initializes all application dependencies
// This is the single place where ALL dependencies are wired together
func initDependencies(db *gorm.DB) {
	//  INITIALIZE REPOSITORIES
	productRepository = repository.NewProductRepository(db)
	userRepository = repository.NewUserRepository(db)
	cartRepository = repository.NewCartRepository(db)

	//  INITIALIZE SERVICES
	userService = services.NewUserService(userRepository)
	productService = services.NewProductService(productRepository)
	cartService := services.NewCartService(cartRepository, productRepository)

	//  INITIALIZE CONTROLLERS
	authController = controller.NewAuthController(userService)
	productController = controller.NewProductController(productService)
	userController = controller.NewUserController(userService)
	cartController = controller.NewCartController(cartService)
}
