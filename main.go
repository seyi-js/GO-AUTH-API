package main

import (
	"fmt"
	"net/http"
	"os"

	"go-auth-api/src/configs"
	"go-auth-api/src/responses"
	"go-auth-api/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	configs.LoadEnv()
	PORT := os.Getenv("PORT")
	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	configs.ConnectDB()

	v1 := router.Group("/api/v1")

	routes.UserROute(v1)

	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.ExampleResponse{Message: "ok"})
	})

	fmt.Println("Server running on port:" + PORT)
	router.Run(":" + PORT)

}
