package controllers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nNottp33/project-go-fiber/src/middleware"
	models "github.com/nNottp33/project-go-fiber/src/models"
	"gorm.io/gorm"
)

func GetMemberProfile(c *fiber.Ctx) error {
	session := c.Locals("session").(*middleware.Sessions)
	var member models.Members
	exec := db.Where("id = ?", session.Id).Last(&member)
	if exec.Error != nil {
		fmt.Println("[ERROR] GetMemberProfile CATCH:", exec.Error)
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":   fiber.StatusNotFound,
				"errors": "Member not found",
			})
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": exec.Error.Error(),
		})
	}

	member.Password = "xxxxx"
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully",
		"data":    member,
	})

}
