package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/services"
)

type CartController struct {
	cartService services.ICartService
}

func NewCartController(cartService services.ICartService) *CartController {
	return &CartController{cartService: cartService}
}

func (cc *CartController) AddToCart(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}

	var req models.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := cc.cartService.AddToCart(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart"})
}
