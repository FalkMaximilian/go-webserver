package routes

import (
	"go-webserver/handlers"
	"go-webserver/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupSetRoutes(app fiber.Router) {
	app.Post("/sets", middleware.JWTProtected(), handlers.CreateSetHandler)
	app.Delete("/sets/:set_id<int>", middleware.JWTProtected(), handlers.DeleteSetHandler)
	app.Post("/sets/:set_id<int>/cards", middleware.JWTProtected(), handlers.AddCardToSetHandler)
}
