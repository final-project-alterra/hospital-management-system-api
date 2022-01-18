package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type env struct {
	PORT       string
	JWT_SECRET string

	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_TIMEZONE string

	TIMEZONE string
}

var ENV env

func LoadENV(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Println("Failed loading '.env' file. Error:", err.Error())
	}

	ENV.PORT = os.Getenv("PORT")
	if ENV.PORT == "" {
		ENV.PORT = "8080"
	}

	ENV.JWT_SECRET = os.Getenv("JWT_SECRET")

	ENV.DB_HOST = os.Getenv("DB_HOST")
	ENV.DB_PORT = os.Getenv("DB_PORT")
	ENV.DB_NAME = os.Getenv("DB_NAME")
	ENV.DB_USER = os.Getenv("DB_USER")
	ENV.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	ENV.DB_TIMEZONE = os.Getenv("DB_TIMEZONE")

	ENV.TIMEZONE = os.Getenv("TIMEZONE")
}
