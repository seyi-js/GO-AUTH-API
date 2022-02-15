package controllers

import (
	"context"
	"fmt"
	"go-auth-api/src/configs"
	"go-auth-api/src/models"
	"go-auth-api/src/responses"
	"go-auth-api/src/utils"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
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

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Info: map[string]interface{}{"info": err}})

			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, responses.Response{Status: http.StatusConflict, Info: map[string]interface{}{"info": "email exists."}})

			return
		}

		var HashedPassword string = utils.HashPassword(user.Password)
		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Email:    user.Email,
			Password: HashedPassword,
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

func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var user models.User
		var foundUser models.User

		defer cancel()

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

		err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Info: map[string]interface{}{"info": "invalid login credentials."}})
			return
		}

		isValid, msg := utils.ComparePassword(foundUser.Password, user.Password)

		defer cancel()

		fmt.Println(msg)

		if !isValid {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Info: map[string]interface{}{"info": msg}})
			return
		}

		token, err := utils.GenerateJWT(foundUser.Id.Hex())

		if err != nil {
			if !isValid {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Info: map[string]interface{}{"info": err.Error()}})
				return
			}
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Info: map[string]interface{}{"info": gin.H{
			"token": token,
		}}})

	}
}
