package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading env file.")
	}
}

func EnvMongoURI() string {
	LoadEnv()
	uri := os.Getenv("MONGOURI")
	return uri
}
