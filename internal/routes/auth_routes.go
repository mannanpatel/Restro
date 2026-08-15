package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/controller"
	"github.com/mannanpatel/Restro/internal/interfaces"
	"github.com/mannanpatel/Restro/internal/repository"
	"github.com/mannanpatel/Restro/internal/service"
)

var (
	userRepository interfaces.IUserRepository = repository.NewUserRepository()
	userService    interfaces.IUserService    = service.NewUserService(userRepository)
	authController interfaces.IAuthController = controller.NewAuthController(userService)
)

func registerAuthRoutes(r *gin.Engine) {
	r.POST("/login", authController.Login)
	r.POST("/register", authController.Register)

	// admin := r.Group("/admin")
	// admin.Use(middleware.JWTMiddleware(), middleware.AdminOnly())
	// admin.PUT("/users/:id/role", authController.UpdateUserRole)

}
