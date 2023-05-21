package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/nNottp33/project-go-fiber/src/middleware"
	models "github.com/nNottp33/project-go-fiber/src/models"
	util "github.com/nNottp33/project-go-fiber/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	AdminUsersBody struct {
		Username string `json:"username" xml:"username" format:"username" validate:"required"`
		Password string `json:"password" xml:"password" format:"password" validate:"required"`
		Role     string `json:"role" xml:"role" format:"role" validate:"required"`
		FullName string `json:"full_name" xml:"full_name" format:"full_name" validate:"required"`
	}

	UpdateAdminUsersBody struct {
		Password  string    `json:"password" xml:"password" format:"password"`
		Role      string    `json:"role" xml:"role" format:"role"`
		FullName  string    `json:"full_name" xml:"full_name" format:"full_name"`
		UpdatedBy string    `json:"updated_by"`
		UpdatedAt time.Time `json:"updated_at"`
		IsActive  bool      `json:"is_active"`
	}
)

func GetAdmins(c *fiber.Ctx) error {
	var admin []models.AdminUsers
	if exec := db.Find(&admin).Order("updated_at DESC"); exec.Error != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": exec.Error.Error(),
		})
	}

	for i := 0; i < len(admin); i++ {
		admin[i].Password = "xxxxxx"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully",
		"data":    admin,
	})
}

func GetAdminById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		fmt.Println("[ERROR] GetAdminById params CATCH:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": err.Error(),
		})
	}

	var admin models.AdminUsers
	exec := db.Where("id = ?", id).Last(&admin)
	admin.Password = "xxxxxx"
	if exec.Error != nil {
		fmt.Println("[ERROR] GetAdminById CATCH:", exec)
		if errors.Is(exec.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":   fiber.StatusNotFound,
				"errors": "Admin not found",
			})
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": exec.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully",
		"data":    admin,
	})
}

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
		fmt.Println("[ERROR] CreateAdmin hash CATCH:", errHash)
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

func EditAdmin(c *fiber.Ctx) error {
	id, errParams := c.ParamsInt("id")
	if errParams != nil {
		fmt.Println("[ERROR] EditAdmin params CATCH:", errParams)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errParams.Error,
		})
	}

	body := new(UpdateAdminUsersBody)
	if errBodyParser := c.BodyParser(body); errBodyParser != nil {
		fmt.Println("[ERROR] EditAdmin BodyParser CATCH:", errBodyParser)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errBodyParser.Error(),
		})
	}

	salt, errHash := bcrypt.GenerateFromPassword([]byte(body.Password), 11)
	if errHash != nil {
		fmt.Println("[ERROR] EditAdmin map model CATCH:", errHash)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": errHash,
		})
	}

	var adminUsers models.AdminUsers
	session := c.Locals("session").(*middleware.Sessions)
	body.UpdatedBy = session.Username
	body.Password = string(salt)
	body.UpdatedAt = time.Now()
	if exec := db.Model(adminUsers).Where("id = ?", id).Updates(&body); exec.Error != nil {
		fmt.Println("[ERROR] EditAdmin exec CATCH:", exec.Error)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": exec.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Update admin successfully",
	})
}

func DeleteAdmin(c *fiber.Ctx) error {
	id, errParams := c.ParamsInt("id")
	if errParams != nil {
		fmt.Println("[ERROR] DeleteAdmin params CATCH:", errParams)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errParams.Error,
		})
	}

	var adminUsers models.AdminUsers
	if exec := db.Delete(&adminUsers, id); exec.Error != nil {
		fmt.Println("[ERROR] DeleteAdmin exec CATCH:", exec.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": exec.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Delete admin successfully",
	})
}
