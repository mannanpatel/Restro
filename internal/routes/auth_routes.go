package routes

import (
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.Engine) {

	r.POST("/login", authController.Login)
	r.POST("/register", authController.Register)

	// admin := r.Group("/admin")
	// admin.Use(middleware.JWTMiddleware(), middleware.AdminOnly())
	// admin.PUT("/users/:id/role", authController.UpdateUserRole)

}
