package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nNottp33/project-go-fiber/src/controllers"
	"github.com/nNottp33/project-go-fiber/src/middleware"
)

func IndexRoute(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"protocol":      c.Protocol(),
			"headers":       c.GetReqHeaders(),
			"request_ip":    c.IP(),
			"request_proxy": c.IPs(),
			"current_time":  time.Now(),
			"message":       "Test",
		})
	})

	auth := app.Group("/auth")
	AuthRoutes(auth)

	common := app.Group("/api")
	common.Use(middleware.AuthenticationMiddleware)
	common.Use(middleware.SessionMiddleware)
	BookRoutes(common)
	AdminRoutes(common)
	MemberRoutes(common)
}

func MemberRoutes(router fiber.Router) {
	router.Get("/member", controllers.GetMemberProfile)
}

func AdminRoutes(router fiber.Router) {
	router.Get("/admins", controllers.GetAdmins)
	router.Get("/admin/:id", controllers.GetAdminById)
	router.Post("/admin", controllers.CreateAdmin)
	router.Put("/admin/:id", controllers.EditAdmin)
	router.Delete("/admin/:id", controllers.DeleteAdmin)
}

func BookRoutes(router fiber.Router) {
	router.Get("/books", controllers.GetBooks)
	router.Get("/book/:id", controllers.GetBookById)
	router.Post("/book", controllers.NewBook)
	router.Put("/book/:id", controllers.EditBook)
	router.Delete("/book/:id", controllers.DeleteBook)
}

func AuthRoutes(router fiber.Router) {
	router.Post("/login", controllers.SignIn)
	router.Patch("/logout", controllers.SignOut)
}
