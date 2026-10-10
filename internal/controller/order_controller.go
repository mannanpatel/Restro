package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/repository"
	"github.com/mannanpatel/Restro/internal/services"
	"gorm.io/gorm"
)

type OrderController struct {
	orderService services.IOrderService
}

func NewOrderController(orderService services.IOrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

type updateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed preparing ready completed cancelled"`
}

func (oc *OrderController) CreateFromCart(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}
	order, err := oc.orderService.CreateFromCart(c.Request.Context(), userID)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (oc *OrderController) ListMine(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}
	orders, err := oc.orderService.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) GetMine(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}
	orderID, err := parseOrderID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	order, err := oc.orderService.GetByID(c.Request.Context(), userID, orderID)
	if err != nil {
		writeOrderError(c, err)
		return
	}
	c.JSON(http.StatusOK, order)
}

func (oc *OrderController) ListAll(c *gin.Context) {
	orders, err := oc.orderService.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) UpdateStatus(c *gin.Context) {
	orderID, err := parseOrderID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	var req updateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := oc.orderService.UpdateStatus(c.Request.Context(), orderID, req.Status); err != nil {
		writeOrderError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func parseOrderID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid order id")
	}
	return uint(id), nil
}

func writeOrderError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrCartEmpty):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
	case errors.Is(err, repository.ErrProductUnavailable):
		c.JSON(http.StatusConflict, gin.H{"error": "A cart product is unavailable"})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order operation failed"})
	}
}
