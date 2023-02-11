package config

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/joho/godotenv"
)

var DbHost string
var DbUser string
var DbPassword string
var DbName string
var DbPort string
var ServerPort string
var JWTSecret string
var kafkaUrl string

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(color.RedString("Error loading .env"))
	}

	DbHost = os.Getenv("DB_HOST")
	DbUser = os.Getenv("POSTGRES_USER")
	DbPassword = os.Getenv("POSTGRES_PASSWORD")
	DbName = os.Getenv("POSTGRES_DB")
	DbPort = os.Getenv("POSTGRES_PORT")
	ServerPort = os.Getenv("SERVER_PORT")
	JWTSecret = os.Getenv("JWT_SECRET")
	kafkaUrl = os.Getenv("KAFKA_URL")
}
