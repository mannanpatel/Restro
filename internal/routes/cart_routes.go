package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/middleware"
)

func registerCartRoutes(r *gin.Engine) {
	cart := r.Group("/cart")
	cart.Use(middleware.JWTMiddleware())

	cart.GET("", cartController.GetCart)
	cart.POST("/add", cartController.AddToCart)
	cart.POST("/items", cartController.AddToCart)
	cart.PATCH("/items/:itemID", cartController.UpdateCartItem)
	cart.DELETE("/items/:itemID", cartController.RemoveCartItem)
	cart.DELETE("", cartController.ClearCart)
}

// Cart item access is scoped to the authenticated user's cart.
