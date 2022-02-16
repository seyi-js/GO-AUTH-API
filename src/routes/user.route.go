package routes

import (
	"go-auth-api/src/controllers"
	"go-auth-api/src/middlewares"

	"github.com/gin-gonic/gin"
)

func UserROute(router *gin.RouterGroup) {
	router.POST("/register", controllers.SignUp())
	router.POST("/login", controllers.SignIn())
	router.GET("/me", middlewares.AuthorizedUser(), controllers.GetUser())
}
