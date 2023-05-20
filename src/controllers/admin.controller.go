package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/nNottp33/project-go-fiber/src/middleware"
	models "github.com/nNottp33/project-go-fiber/src/models"
	util "github.com/nNottp33/project-go-fiber/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

type (
	AdminUsersBody struct {
		Username string `json:"username" xml:"username" format:"username" validate:"required"`
		Password string `json:"password" xml:"password" format:"password" validate:"required"`
		Role     string `json:"role" xml:"role" format:"role" validate:"required"`
		FullName string `json:"full_name" xml:"full_name" format:"full_name" validate:"required"`
	}
)

func CreateAdmin(c *fiber.Ctx) error {
	body := new(AdminUsersBody)
	if err := c.BodyParser(body); err != nil {
		fmt.Println("[ERROR] CreateAdmin CATCH:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": err.Error(),
		})
	}

	if errValidate := util.ValidateBody(body); errValidate != nil {
		fmt.Println("[ERROR] CreateAdmin validate payload CATCH:", errValidate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errValidate,
		})
	}

	var adminUsers models.AdminUsers
	if errCopy := copier.Copy(&adminUsers, &body); errCopy != nil {
		fmt.Println("[ERROR] CreateAdmin map model CATCH:", errCopy)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": errCopy,
		})
	}

	salt, errHash := bcrypt.GenerateFromPassword([]byte(body.Password), 11)
	if errHash != nil {
		fmt.Println("[ERROR] CreateAdmin map model CATCH:", errHash)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": errHash,
		})
	}

	session := c.Locals("session").(*middleware.Sessions)
	adminUsers.Password = string(salt)
	adminUsers.CreatedBy = session.Username
	adminUsers.CreatedAt = time.Now()
	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "username"}}, DoNothing: true}).Create(&adminUsers)
	if result.Error != nil {
		fmt.Println("[ERROR] CreateAdmin Exec query CATCH:", result.Error)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": result.Error,
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"code":   fiber.StatusConflict,
			"errors": "Username already exists",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Created admin successfully",
	})
}
