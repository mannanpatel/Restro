package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/middleware"
)

func registerProductRoutes(r *gin.Engine) {
	// ==================== PUBLIC PRODUCT ROUTES ====================
	// No authentication required for read operations
	productPublic := r.Group("/products")
	{
		productPublic.GET("", productController.GetAllProducts)
		productPublic.GET("/:id", productController.GetProductById)
	}

	// ==================== ADMIN ONLY PRODUCT ROUTES ====================
	// Admin access required for create/update/delete operations
	productPrivate := r.Group("/admin/products")
	productPrivate.Use(middleware.JWTMiddleware())
	productPrivate.Use(middleware.AdminOnly())
	{
		productPrivate.POST("", productController.CreateProduct)
		productPrivate.PATCH("/:id", productController.UpdateProduct)
		productPrivate.DELETE("/:id", productController.DeleteProduct)
	}
}

// GET    /products
// GET    /products/:id
// POST   /products
// PATCH  /products/:id
// DELETE /products/:id

// GET    /admin/products
// GET    /admin/products/:id
// POST   /admin/products
// PATCH  /admin/products/:id
// DELETE /admin/products/:id
