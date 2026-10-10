package routes

import (
	"github.com/mannanpatel/Restro/config"
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
	orderController   *controller.OrderController
	paymentController *controller.PaymentController
)

// SERVICES
var (
	userService    services.IUserService
	productService services.IProductService
	cartService    services.ICartService
	orderService   services.IOrderService
	paymentService services.IPaymentService
)

// REPOSITORIES
var (
	userRepository    repository.IUserRepository
	productRepository repository.IProductRepository
	cartRepository    repository.ICartRepository
	orderRepository   repository.IOrderRepository
	paymentRepository repository.IPaymentRepository
)

// initDependencies initializes all application dependencies
// This is the single place where ALL dependencies are wired together
func initDependencies(db *gorm.DB) {
	//  INITIALIZE REPOSITORIES
	productRepository = repository.NewProductRepository(db)
	userRepository = repository.NewUserRepository(db)
	cartRepository = repository.NewCartRepository(db)
	orderRepository = repository.NewOrderRepository(db)
	paymentRepository = repository.NewPaymentRepository(db)

	//  INITIALIZE SERVICES
	userService = services.NewUserService(userRepository)
	productService = services.NewProductService(productRepository)
	cartService = services.NewCartService(cartRepository, productRepository)
	orderService = services.NewOrderService(orderRepository)
	paymentGateway := services.NewRazorpayClient(config.AppConfig.Razorpay.KeyID, config.AppConfig.Razorpay.KeySecret)
	paymentService = services.NewPaymentService(paymentRepository, orderService, paymentGateway, config.AppConfig.Razorpay.KeySecret)

	//  INITIALIZE CONTROLLERS
	authController = controller.NewAuthController(userService)
	productController = controller.NewProductController(productService)
	userController = controller.NewUserController(userService)
	cartController = controller.NewCartController(cartService)
	orderController = controller.NewOrderController(orderService)
	paymentController = controller.NewPaymentController(paymentService)
}
