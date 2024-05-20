package routes

import (
	handlers "go-webserver/handlers/api"
	"go-webserver/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAPIRoutes(app *fiber.App) {
	api := app.Group("/api", middleware.JWTProtected())
	api.Get("/protected-test", handlers.ProtectedTestHandler)
}
