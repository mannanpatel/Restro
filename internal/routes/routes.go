package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	initDependencies()

	registerAuthRoutes(r)
	registerProductRoutes(r)
	registerUserRoutes(r)

	// registerCartRoutes(r)
	// registerOrderRoutes(r)
	// registerPaymentRoutes(r)
}
