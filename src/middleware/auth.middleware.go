package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	configs "github.com/nNottp33/project-go-fiber/src/configs"
)

type Sessions struct {
	Id       int
	Username string
	Role     string
}

var db = configs.Db

func AuthenticationMiddleware(c *fiber.Ctx) error {
	token, errCheckHeaders := checkHeaders(c)
	if errCheckHeaders != nil {
		fmt.Println("[ERROR] AuthenticationMiddleware checkHeaders CATCH:", errCheckHeaders)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":  fiber.StatusUnauthorized,
			"error": errCheckHeaders.Error(),
		})
	}

	result, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(configs.JWT_SECRET), nil
	})

	if err != nil {
		fmt.Println("[ERROR] AuthenticationMiddleware Verify token CATCH:", token)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":  fiber.StatusUnauthorized,
			"error": []string{"Unauthorized", err.Error()},
		})
	}

	claims, ok := result.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("[ERROR] AuthenticationMiddleware claims CATCH:", ok)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":  fiber.StatusUnauthorized,
			"error": "Unauthorized",
		})
	}

	c.Locals("session", &Sessions{
		Id:       int(claims["id"].(float64)),
		Username: claims["username"].(string),
		Role:     claims["role"].(string),
	})

	return c.Next()
}

func SessionMiddleware(c *fiber.Ctx) error {
	token, err := checkHeaders(c)
	if err != nil {
		fmt.Println("[ERROR] SessionMiddleware checkHeaders CATCH:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":  fiber.StatusUnauthorized,
			"error": err.Error(),
		})
	}

	session := c.Locals("session").(*Sessions)
	var table = "admin_users"
	if session.Role == "member" {
		table = "member"
	}

	var isExists int64
	errors := db.Table(table).Where("id = ?", session.Id).Where("session_token = ?", token).Count(&isExists)
	if errors.Error != nil || isExists == 0 {
		fmt.Println("[ERROR] SessionMiddleware check token CATCH:", errors.Error, "isExists:", isExists)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":  fiber.StatusUnauthorized,
			"error": "Session expired!",
		})
	}

	return c.Next()
}

func checkHeaders(c *fiber.Ctx) (string, error) {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token == "" {
		return "", errors.New("invalid token")
	}

	return token, nil
}
