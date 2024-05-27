package routes

import (
	"go-webserver/handlers"
	"go-webserver/logger"
	"go-webserver/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupSetRoutes(app fiber.Router) {
	logger.Log.Info("setting up routes for set crud operations")
	app.Get("/sets", middleware.JWTProtected(), handlers.GetSetsHandler)
	app.Get("/sets/:set_id<int>", middleware.JWTProtected(), handlers.GetSetHandler)
	app.Post("/sets", middleware.JWTProtected(), handlers.CreateSetHandler)
	app.Delete("/sets/:set_id<int>", middleware.JWTProtected(), handlers.DeleteSetHandler)
	app.Post("/sets/:set_id<int>/cards", middleware.JWTProtected(), handlers.AddCardToSetHandler)
	app.Delete("/sets/:set_id<int>/cards/:card_id<int>", middleware.JWTProtected(), handlers.RemoveCardFromSetHandler)
}
