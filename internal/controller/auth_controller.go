package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mannanpatel/Restro/internal/interfaces"
	"github.com/mannanpatel/Restro/internal/models"
)

type authController struct {
	service interfaces.IUserService
}

func NewAuthController(service interfaces.IUserService) interfaces.IAuthController {
	return &authController{service: service}
}

func (ac *authController) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := ac.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (ac *authController) Register(c *gin.Context) {
	var user models.UserDetails
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.service.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// func (ac *authController) Register(c *gin.Context) {
// 	var req models.UserRegisterRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := ac.service.RegisterUser(c.Request.Context(), req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
// }

// func (ac *authController) UpdateUserRole(c *gin.Context) {
// 	// TODO: implement after role management is defined
// 	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
// }

// claims := jwt.MapClaims{
// 		"user_id": user.ID,
// 		"email":   user.Email,
// 		"role":    user.Role,
// 		"exp":     time.Now().Add(24 * time.Hour).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signed, err := token.SignedString([]byte(ac.jwtSecret))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.UserLoginResponse{
// 		UserId: user.ID,
// 		Email:  user.Email,
// 		Role:   user.Role,
// 		Token:  signed,
// 	})
