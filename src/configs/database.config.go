package configs

import (
	"fmt"
	"time"

	"github.com/nNottp33/project-go-fiber/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok", Database.Host, Database.User, Database.Pass, Database.Name, Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("[ERROR] Failed to connect to database")
	}

	errMigrate := db.AutoMigrate(&models.AdminUsers{}, &models.Members{}, &models.Book{})
	if errMigrate != nil {
		panic("[ERROR] Can not migrate tables")
	}

	sqlDB, sqlDBErr := db.DB()
	if sqlDBErr != nil {
		panic(sqlDBErr)
	}

	sqlDB.SetMaxIdleConns(Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Minute * 15)
	Db = db
}
