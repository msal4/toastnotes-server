package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/msal4/toastnotes/middleware"
	"gorm.io/gorm"
)

const (
	// API is the v1 api group.
	API = "/api/v1"
	// APIRegister is the user registeration endpoint.
	APIRegister = "/register"
	// APILogin is the user signin endpoint.
	APILogin = "/login"
	// APIRefresh is the user tokens refresh endpoint.
	APIRefresh = "/refresh"
	// APILogout is the logout endpoint.
	APILogout = "/logout"
	// APIChangePassword is the authenticated user endpoint for changing their password.
	APIChangePassword = "/change_password"
	// APIMe is the user profile endpoint.
	APIMe = "/me"

	// APINote is the user notes api group.
	APINote = "/notes"
)

// SetupRouter sets up the app routes.
func SetupRouter(db *gorm.DB) *gin.Engine {
	// router
	router := gin.New()

	// middleware
	router.Use(gin.Recovery(), middleware.CORS())

	// controllers
	userController := NewUserController(db)
	noteController := NewNoteController(db)

	v1 := router.Group(API)
	{
		v1.POST(APIRegister, userController.Register)
		v1.POST(APILogin, userController.Login)
		v1.POST(APIRefresh, userController.RefreshTokens)

		authenticated := v1.Group("/", middleware.JWTAuth())
		{
			// user
			authenticated.GET(APIMe, userController.Me)
			authenticated.POST(APIChangePassword, userController.ChangePassword)
			authenticated.DELETE(APILogout, userController.Logout)

			// note
			authenticated.GET(APINote, noteController.List)
			authenticated.POST(APINote, noteController.Create)
			authenticated.GET(APINote+"/:id", noteController.Retrieve)
			authenticated.PUT(APINote+"/:id", noteController.Update)
			authenticated.DELETE(APINote+"/:id", noteController.Delete)
		}
	}

	return router
}
