package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	configs "github.com/nNottp33/project-go-fiber/src/configs"
	util "github.com/nNottp33/project-go-fiber/src/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	Login struct {
		Username string `json:"username" xml:"username" form:"username" validate:"required"`
		Password string `json:"password" xml:"password" form:"password" validate:"required"`
		Source   string `json:"source" xml:"source" form:"source" validate:"required,oneof=admin member"`
	}

	Logout struct {
		Id     int8   `json:"id" xml:"id" form:"id"`
		Source string `json:"source" xml:"source" form:"source"`
	}

	Profile struct {
		Id           int8   `json:"id"`
		Username     string `json:"username"`
		Password     string `json:"-"`
		FullName     string `json:"full_name"`
		Role         string `json:"role"`
		SessionToken string `json:"session_token"`
		IsLogin      bool   `json:"is_login"`
	}
)

var profile = Profile{}

func SignIn(c *fiber.Ctx) error {
	body := new(Login)
	if err := c.BodyParser(body); err != nil {
		fmt.Println("[ERROR] SignIn BodyParser CATCH:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": err.Error(),
		})
	}

	if errValidate := util.ValidateBody(body); errValidate != nil {
		fmt.Println("[ERROR] SignIn validate payload CATCH:", errValidate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errValidate,
		})
	}

	table := getModel(body.Source)
	if checkUser := db.Table(table).Where("username = ?", body.Username).First(&profile); checkUser.Error != nil {
		fmt.Println("[ERROR] SignIn checkUser CATCH:", checkUser.Error)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": "User not found",
		})
	}

	hash := []byte(profile.Password)
	password := []byte(body.Password)
	if isMatch := bcrypt.CompareHashAndPassword(hash, password); isMatch != nil {
		fmt.Println("[ERROR] SignIn errComparePassword CATCH:", isMatch)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":   fiber.StatusUnauthorized,
			"errors": "Password incorrect!",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        profile.Id,
		"username":  profile.Username,
		"full_name": profile.FullName,
		"role":      profile.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, errJwt := token.SignedString([]byte(configs.JWT_SECRET))
	if errJwt != nil {
		fmt.Println("[ERROR] SignIn errJwt CATCH:", errJwt)
		return c.Status(fiber.StatusFailedDependency).JSON(fiber.Map{
			"code":   fiber.StatusFailedDependency,
			"errors": errJwt,
		})
	}

	profile.SessionToken = tokenString
	profile.IsLogin = true
	isUpdate, errUpdateToken := updateToken(table, profile)
	if !isUpdate {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": errUpdateToken,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Sign In successfully",
		"data":    profile,
	})
}

func SignOut(c *fiber.Ctx) error {
	body := new(Logout)
	if err := c.BodyParser(body); err != nil {
		fmt.Println("[ERROR] SignOut BodyParser CATCH:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": err.Error(),
		})
	}

	if errValidate := util.ValidateBody(body); errValidate != nil {
		fmt.Println("[ERROR] SignOut validate payload CATCH:", errValidate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":   fiber.StatusBadRequest,
			"errors": errValidate,
		})
	}

	table := getModel(body.Source)
	profile.SessionToken = ""
	isUpdate, errUpdateToken := updateToken(table, profile)
	if !isUpdate {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":   fiber.StatusUnprocessableEntity,
			"errors": errUpdateToken,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "SignOut successfully",
	})
}

func getModel(source string) string {
	var table = "admin_users"
	if source == "member" {
		table = "members"
	}

	return table
}

func updateToken(table string, profile Profile) (bool, *gorm.DB) {
	errUpdateToken := db.Table(table).Where("id = ?", profile.Id).Update("session_token", profile.SessionToken)
	if errUpdateToken.Error != nil {
		fmt.Println("[ERROR] SignIn errUpdateToken CATCH:", errUpdateToken.Error)
		return false, errUpdateToken
	}

	return true, nil
}
