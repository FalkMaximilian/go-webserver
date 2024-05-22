package routes

import (
	handlers "go-webserver/handlers/api"
	"go-webserver/logger"
	"go-webserver/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAPIRoutes(app *fiber.App) {
	logger.Info("Setting up protected routes for /api")
	api := app.Group("/api", middleware.JWTProtected())
	api.Get("/protected-test", handlers.ProtectedTestHandler)
	api.Post("/sets", handlers.CreateSet)
}
