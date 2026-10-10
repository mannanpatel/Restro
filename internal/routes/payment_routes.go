package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/middleware"
)

func registerPaymentRoutes(r *gin.Engine) {
	orders := r.Group("/orders", middleware.JWTMiddleware())
	orders.POST("/:id/payment", paymentController.CreateOrder)

	payments := r.Group("/payments", middleware.JWTMiddleware())
	payments.POST("/verify", paymentController.Verify)
}
