package routes

import (
	"github.com/mannanpatel/Restro/internal/controller"
	"github.com/mannanpatel/Restro/internal/repository"
	"github.com/mannanpatel/Restro/internal/services"
)

// CONTROLLERS
var (
	authController    *controller.AuthController
	productController *controller.ProductController
	userController    *controller.UserController
)

// SERVICES
var (
	userService    services.IUserService
	productService services.IProductService
)

// REPOSITORIES
var (
	userRepository    repository.IUserRepository
	productRepository repository.IProductRepository
)

// initDependencies initializes all application dependencies
// This is the single place where ALL dependencies are wired together
func initDependencies() {
	//  INITIALIZE REPOSITORIES
	productRepository = repository.NewProductRepository()
	userRepository = repository.NewUserRepository()

	//  INITIALIZE SERVICES
	userService = services.NewUserService(userRepository)
	productService = services.NewProductService(productRepository)

	//  INITIALIZE CONTROLLERS
	authController = controller.NewAuthController(userService)
	productController = controller.NewProductController(productService)
	userController = controller.NewUserController(userService)
}
