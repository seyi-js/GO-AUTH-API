package routes

import (
	"go-auth-api/src/controllers"

	"github.com/gin-gonic/gin"
)

func UserROute(router *gin.RouterGroup) {
	router.POST("/register", controllers.SignUp())
}
