package main

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nNottp33/project-go-fiber/src/configs"
	"github.com/nNottp33/project-go-fiber/src/routes"
)

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
	}))

	// get routes
	routes.IndexRoute(app)

	err := app.Listen("localhost:" + configs.Port)
	if err != nil {
		panic(err)
	}
}
