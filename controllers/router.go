package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/msal4/toastnotes/middleware"
	"gorm.io/gorm"
)

// SetupRouter sets up the app routes.
func SetupRouter(db *gorm.DB) *gin.Engine {
	// router
	router := gin.Default()

	router.Use(middleware.CORS())

	// controllers
	userController := NewUserController(db)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", userController.Register)
		v1.POST("/login", userController.Login)
		v1.POST("/refresh", userController.RefreshTokens)

		authenticated := v1.Group("/", middleware.JWTAuth())
		{
			authenticated.GET("/me", userController.Me)
			authenticated.POST("/change_password", userController.ChangePassword)
		}
	}

	return router
}
