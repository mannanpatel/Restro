package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/repository"
	"github.com/mannanpatel/Restro/internal/services"
	"gorm.io/gorm"
)

type PaymentController struct {
	paymentService services.IPaymentService
}

func NewPaymentController(paymentService services.IPaymentService) *PaymentController {
	return &PaymentController{paymentService: paymentService}
}

func (pc *PaymentController) CreateOrder(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}
	orderIDValue, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderIDValue == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	response, err := pc.paymentService.CreateOrder(c.Request.Context(), userID, uint(orderIDValue))
	if err != nil {
		writePaymentError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (pc *PaymentController) Verify(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user identity"})
		return
	}
	var req models.VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := pc.paymentService.VerifyPayment(c.Request.Context(), userID, req); err != nil {
		writePaymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment verified", "status": "paid"})
}

func writePaymentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrPaymentNotConfigured):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Payment service is not configured"})
	case errors.Is(err, services.ErrInvalidPaymentSignature), errors.Is(err, services.ErrProviderPaymentMismatch):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment verification failed"})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment or order not found"})
	case errors.Is(err, repository.ErrPaymentConflict), errors.Is(err, repository.ErrOrderNotPayable):
		c.JSON(http.StatusConflict, gin.H{"error": "Order cannot accept this payment"})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"error": "Payment provider operation failed"})
	}
}
