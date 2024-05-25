package routes

import (
	"go-webserver/handlers"
	"go-webserver/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app fiber.Router) {
	app.Post("users/register", middleware.RedirectIfAuthenticated(), handlers.RegisterUserHandler)
	app.Post("users/login", middleware.RedirectIfAuthenticated(), handlers.LoginHandler)
}
