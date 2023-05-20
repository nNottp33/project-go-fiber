package configs

import (
	"time"

	"github.com/nNottp33/project-go-fiber/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(postgres.Open(Database.URL), &gorm.Config{})
	if err != nil {
		panic("[ERROR] Failed to connect to database")
	}

	db.AutoMigrate(&models.AdminUsers{}, &models.Members{}, &models.Book{})

	sqlDB, sqlDBErr := db.DB()
	if sqlDBErr != nil {
		panic(sqlDBErr)
	}

	sqlDB.SetMaxIdleConns(Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Minute * 15)
	Db = db
}
