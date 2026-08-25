package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/middleware"
)

func registerUserRoutes(r *gin.Engine) {
	// Authenticated users
	users := r.Group("/users", middleware.JWTMiddleware())
	{
		users.GET("/:id", userController.GetUserById)
		users.PATCH("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	// Admin only
	adminUsers := r.Group("/admin/users", middleware.JWTMiddleware(), middleware.AdminOnly())
	{
		adminUsers.GET("", userController.ListUsers)
		adminUsers.PUT("/:id/role", userController.UpdateUserRole)
	}
}
