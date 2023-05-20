package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConnection struct {
	URL          string
	MaxIdleConns int
	MaxOpenConns int
}

var (
	PORT       string
	JWT_SECRET string
	Database   DatabaseConnection
)

func LoadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT = os.Getenv("PORT")
	JWT_SECRET = os.Getenv("JWT_SECRET")
	Database.URL = os.Getenv("DB_URL")

	maxIdleCon, errMaxIdleCon := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	if errMaxIdleCon != nil {
		panic("[ERROR] env 'DB_MAX_IDLE_CONNECTIONS' invalid or not integer value")
	}

	maxOpenCon, errMaxOpenCon := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	if errMaxOpenCon != nil {
		panic("[ERROR] env 'DB_MAX_CONNECTIONS' invalid or not integer value")
	}

	Database.MaxIdleConns = maxIdleCon
	Database.MaxOpenConns = maxOpenCon
}
