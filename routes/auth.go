package routes

import (
	"go-webserver/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	app.Post("/register", handlers.RegisterUserHandler)
	app.Post("/login", handlers.LoginHandler)
}
