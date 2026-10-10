package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/middleware"
)

func registerOrderRoutes(r *gin.Engine) {
	orders := r.Group("/orders", middleware.JWTMiddleware())
	{
		orders.POST("", orderController.CreateFromCart)
		orders.GET("", orderController.ListMine)
		orders.GET("/:id", orderController.GetMine)
	}

	adminOrders := r.Group("/admin/orders", middleware.JWTMiddleware(), middleware.AdminOnly())
	{
		adminOrders.GET("", orderController.ListAll)
		adminOrders.PATCH("/:id/status", orderController.UpdateStatus)
	}
}
