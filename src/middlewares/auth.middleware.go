package middlewares

import (
	"fmt"
	"go-auth-api/src/responses"
	"go-auth-api/src/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, responses.Response{Status: http.StatusUnauthorized, Info: map[string]interface{}{"info": "No token, authorization denied."}})

			c.Abort()
			return
		}
		clientToken := strings.ReplaceAll(authHeader[len(BEARER_SCHEMA):], " ", "")

		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, responses.Response{Status: http.StatusUnauthorized, Info: map[string]interface{}{"info": "No token, authorization denied."}})

			c.Abort()
			return
		}

		claims, err := utils.VerifyJWT(clientToken)

		if err != "" {
			c.JSON(http.StatusUnauthorized, responses.Response{Status: http.StatusUnauthorized, Info: map[string]interface{}{"info": "invalid token, authorization denied."}})

			fmt.Print(err)

			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
