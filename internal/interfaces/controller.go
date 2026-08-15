package interfaces

import "github.com/gin-gonic/gin"

type IAuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	// UpdateUserRole(c *gin.Context)
}
