package routes

import (
	"go-webserver/handlers"
	"go-webserver/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupCardRoutes(app fiber.Router) {
	app.Get("/cards", middleware.JWTProtected(), handlers.GetCardsHandler)
	app.Get("/cards/:card_id<int>", middleware.JWTProtected(), handlers.GetCardHandler)
	app.Post("/cards", middleware.JWTProtected(), handlers.CreateCardHandler)
	app.Delete("cards/:card_id<int>", middleware.JWTProtected(), handlers.DeleteCardHandler)
}
