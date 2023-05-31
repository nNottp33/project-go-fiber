package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConnection struct {
	Host         string
	Port         int
	User         string
	Pass         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
}

var (
	Port      string
	JwtSecret string
	Database  DatabaseConnection
)

func LoadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Port = os.Getenv("PORT")
	JwtSecret = os.Getenv("JWT_SECRET")
	Database.Host = os.Getenv("DB_HOST")

	dbPort, errDbPort := strconv.Atoi(os.Getenv("DB_PORT"))
	if errDbPort != nil {
		panic("[ERROR] env 'DB_PORT' invalid or not integer value")
	}

	Database.Port = dbPort
	Database.User = os.Getenv("DB_USER")
	Database.Pass = os.Getenv("DB_PASS")
	Database.Name = os.Getenv("DB_NAME")

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
