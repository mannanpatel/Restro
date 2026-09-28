package routes

import (
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
	}

	// admin := r.Group("/admin")
	// admin.Use(middleware.JWTMiddleware(), middleware.AdminOnly())
	// admin.PUT("/users/:id/role", authController.UpdateUserRole)

}

// POST   /auth/register
// POST   /auth/login
// POST   /auth/refresh
// POST   /auth/logout
