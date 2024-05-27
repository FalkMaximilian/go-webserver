package routes

import (
	"go-webserver/handlers"
	"go-webserver/logger"
	"go-webserver/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app fiber.Router) {
	logger.Log.Info("setup routes for user related actions")
	app.Post("/users/register", middleware.RedirectIfAuthenticated(), handlers.RegisterUserHandler)
	app.Post("/users/login", middleware.RedirectIfAuthenticated(), handlers.LoginHandler)
	app.Delete("/users/delete", middleware.JWTProtected(), handlers.DeleteUserHandler)
}
