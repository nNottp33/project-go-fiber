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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	BookBody struct {
		BookNameTh string  `json:"book_name_th" xml:"book_name_th" form:"book_name_th" validate:"required"`
		BookNameEn string  `json:"book_name_en" xml:"book_name_en" form:"book_name_en" validate:"required"`
		Price      float32 `json:"price" xml:"price" form:"price" validate:"required,number"`
		ImageUrl   string  `json:"image_url" xml:"image_url" form:"image_url"`
		Author     string  `json:"author" xml:"author" form:"author"`
		Publisher  string  `json:"publisher" xml:"publisher" form:"publisher"`
		BookCode   string  `json:"book_code" xml:"book_code" form:"book_code" validate:"required"`
	}

	UpdateBookBody struct {
		BookNameTh string    `json:"book_name_th" xml:"book_name_th" form:"book_name_th"`
		BookNameEn string    `json:"book_name_en" xml:"book_name_en" form:"book_name_en"`
		Price      float32   `json:"price" xml:"price" form:"price" validate:"number"`
		ImageUrl   string    `json:"image_url" xml:"image_url" form:"image_url"`
		Author     string    `json:"author" xml:"author" form:"author"`
		Publisher  string    `json:"publisher" xml:"publisher" form:"publisher"`
		BookCode   string    `json:"book_code" xml:"book_code" form:"book_code"`
		UpdatedBy  string    `json:"updated_by"`
		UpdatedAt  time.Time `json:"updated_at"`
		IsActive   bool      `json:"is_active" xml:"is_active" form:"is_active"`
	}
)

func GetBooks(c *fiber.Ctx) error {
	var books []models.Book
	if result := db.Find(&books).Order("updated_at DESC, id DESC"); result.Error != nil {
		fmt.Println("[ERROR] GetBooks Exec query CATCH:", result.Error)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": result,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully",
		"data":    books,
	})
}

func GetBookById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		fmt.Println("[ERROR] getBookById CATCH:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   400,
			"errors": err.Error(),
		})
	}

	var book models.Book
	if result := db.Where("id = ?", id).Last(&book); result.Error != nil {
		fmt.Println("[ERROR] getBookById Exec query CATCH:", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			book = models.Book{}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":   fiber.StatusNotFound,
				"errors": "Book not found",
			})
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": result.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully",
		"data":    book,
	})
}

func NewBook(c *fiber.Ctx) error {
	body := new(BookBody)
	if err := c.BodyParser(body); err != nil {
		fmt.Println("[ERROR] NewBook CATCH:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": err.Error(),
		})
	}

	if errValidate := util.ValidateBody(body); errValidate != nil {
		fmt.Println("[ERROR] NewBook validate payload CATCH:", errValidate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errValidate,
		})
	}

	var book models.Book
	// map body into models
	if errCopy := copier.Copy(&book, &body); errCopy != nil {
		fmt.Println("[ERROR] NewBook map model CATCH:", errCopy)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": errCopy,
		})
	}

	book.CreatedBy = "systemadmin"
	book.UpdatedBy = "systemadmin"
	book.CreatedAt = time.Now()
	result := db.Clauses(clause.OnConflict{
		// conflict column
		Columns: []clause.Column{{Name: "book_code"}},
		// do update set column..
		// be like DO UPDATE SET book_name_th=EXCLUDED.book_name_th
		DoUpdates: clause.AssignmentColumns([]string{"book_name_th", "book_name_en", "price", "image_url", "author", "publisher", "updated_at"}),
	}).Create(&book)

	if result.Error != nil {
		fmt.Println("[ERROR] NewBook Exec query CATCH:", result.Error)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": result.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Save book successfully",
	})
}

func EditBook(c *fiber.Ctx) error {
	id, errParams := c.ParamsInt("id")
	if errParams != nil {
		fmt.Println("[ERROR] EditBook CATCH:", errParams)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errParams.Error(),
		})
	}

	body := new(UpdateBookBody)
	if err := c.BodyParser(body); err != nil {
		fmt.Println("[ERROR] EditBook CATCH:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": err.Error(),
		})
	}

	var book models.Book
	session := c.Locals("session").(*middleware.Sessions)
	body.UpdatedBy = session.Username
	body.UpdatedAt = time.Now()
	if exec := db.Model(book).Where("id = ?", id).Updates(&body); exec.Error != nil {
		fmt.Println("[ERROR] EditBook exec CATCH:", exec.Error)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": exec.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Update book successfully",
	})
}

func DeleteBook(c *fiber.Ctx) error {
	id, errParams := c.ParamsInt("id")
	if errParams != nil {
		fmt.Println("[ERROR] EditBook CATCH:", errParams)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errParams.Error(),
		})
	}

	var book models.Book
	if exec := db.Delete(&book, id); exec.Error != nil {
		fmt.Println("[ERROR] DeleteBook exec CATCH:", exec.Error)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": exec.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Delete book successfully",
	})
}
