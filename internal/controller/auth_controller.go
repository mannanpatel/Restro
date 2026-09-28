package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/common"
	"github.com/mannanpatel/Restro/internal/models"
	"github.com/mannanpatel/Restro/internal/services"
)

type AuthController struct {
	service services.IUserService
}

func NewAuthController(service services.IUserService) *AuthController {
	return &AuthController{service: service}
}

// Register registers a new user
func (ac *AuthController) Register(c *gin.Context) {
	var req models.UserRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ApiErrorResponse("Invalid request body, Error: "+err.Error(), nil))
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Address:   req.Address,
		Type:      req.Type,
		Age:       req.Age,
	}
	err := ac.service.Register(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiErrorResponse("Registration failed, Error: "+err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registered successfully",
	})
}

// Login authenticates a user
func (ac *AuthController) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := ac.service.Login(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
