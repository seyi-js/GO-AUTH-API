package utils

import (
	"log"
	"os"
	"time"

	"go-auth-api/src/configs"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	UserId string
	jwt.StandardClaims
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func ComparePassword(password string, userPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "invalid credentials provided."
		check = false
	}

	return check, msg
}

func GenerateJWT(user_id string) (jwt_token string, err error) {

	configs.LoadEnv()
	var SECRET_KEY string = os.Getenv("JWT_KEY")

	claims := &SignedDetails{
		UserId: user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}

func VerifyJWT(token string) (claims *SignedDetails, msg string) {
	configs.LoadEnv()
	var SECRET_KEY string = os.Getenv("JWT_KEY")
	signed_token, err := jwt.ParseWithClaims(token, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()

		return
	}

	claims, ok := signed_token.Claims.(*SignedDetails)

	if !ok {
		msg = "invalid token"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "expired token"

		return
	}

	return claims, msg
}
