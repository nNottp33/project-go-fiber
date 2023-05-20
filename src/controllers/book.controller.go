package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
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
)

func GetBooks(c *fiber.Ctx) error {
	var books []models.Book
	result := db.Find(&books)
	if result.Error != nil {
		fmt.Println("[ERROR] GetBooks Exec query CATCH:", result)
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
	result := db.First(&book, id)

	if result.Error != nil {
		fmt.Println("[ERROR] getBookById Exec query CATCH:", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			book = models.Book{}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    fiber.StatusNotFound,
				"message": "Book not found",
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
		"message": "Save book data successfully",
	})
}
