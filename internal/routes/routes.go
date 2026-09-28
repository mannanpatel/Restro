package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	initDependencies(db)

	registerAuthRoutes(r)
	registerProductRoutes(r)
	registerUserRoutes(r)
	registerCartRoutes(r)

	// registerOrderRoutes(r)
	// registerPaymentRoutes(r)
}
