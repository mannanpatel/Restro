package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/services"
	"gorm.io/gorm"
)

type CartController struct {
	cartService services.ICartService
}

func NewCartController(cartService services.ICartService) *CartController {
	return &CartController{cartService: cartService}
}

func (cc *CartController) AddToCart(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
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

func (cc *CartController) GetCart(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}

	cart, err := cc.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (cc *CartController) UpdateCartItem(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}
	itemID, err := strconv.Atoi(c.Param("itemID"))
	if err != nil || itemID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}
	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := cc.cartService.UpdateCartItemQuantity(c.Request.Context(), userID, itemID, req.Quantity); err != nil {
		writeCartError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated"})
}

func (cc *CartController) RemoveCartItem(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}
	itemID, err := strconv.Atoi(c.Param("itemID"))
	if err != nil || itemID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}
	if err := cc.cartService.RemoveCartItem(c.Request.Context(), userID, itemID); err != nil {
		writeCartError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (cc *CartController) ClearCart(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}
	if err := cc.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}
	c.Status(http.StatusNoContent)
}

func authenticatedUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int)
	return id, ok && id > 0
}

func writeCartError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
}
