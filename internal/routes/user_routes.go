package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/middleware"
)

func registerUserRoutes(r *gin.Engine) {
	// Authenticated users
	users := r.Group("/users", middleware.JWTMiddleware())
	{
		users.GET("/me", userController.GetUserById)
		users.PATCH("/me", userController.UpdateUser)
		users.DELETE("/me", userController.DeleteUser)
	}

	// Admin only
	adminUsers := r.Group("/admin/users", middleware.JWTMiddleware(), middleware.AdminOnly())
	{
		adminUsers.GET("", userController.ListUsers)
		adminUsers.PUT("/:id/role", userController.UpdateUserRole)

		// GET    /admin/users
		// GET    /admin/users/:id
		// PATCH  /admin/users/:id
		// DELETE /admin/users/:id
		// PUT    /admin/users/:id/role
	}
}
