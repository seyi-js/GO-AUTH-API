package controllers

import (
	"context"
	"go-auth-api/src/configs"
	"go-auth-api/src/models"
	"go-auth-api/src/responses"
	"go-auth-api/src/utils"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "users")

var validate = validator.New()

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var user models.User

		defer cancel()

		//validate the request body using the decalred models, it checks the type string int e.t.c
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Info: map[string]interface{}{"info": err.Error()}})
			return
		}

		validationErr := validate.Struct(&user)

		if validationErr != nil {

			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Info: map[string]interface{}{"info": utils.ValidatorErrorHandler(validationErr)}})
			return
		}

		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Email:    user.Email,
			Password: user.Password,
		}

		result, err := userCollection.InsertOne(ctx, newUser)

		_ = result
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Info: map[string]interface{}{"info": err}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Info: map[string]interface{}{"info": "User registered successfully."}})

	}
}
