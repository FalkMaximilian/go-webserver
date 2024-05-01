package api

import (
  "github.com/gofiber/fiber/v2"
  "go-webserver/logger"
)

func RegisterHandlers(router fiber.Router) {
  logger.Logger.Info("Setting up API endpoints.")
  router.Get("/test1", handleEndpoint1)
  router.Get("/test2", handleEndpoint2)
}

func handleEndpoint1(c *fiber.Ctx) error {
    // Your handler logic for endpoint1
    return c.SendString("Handler for test1")
}

func handleEndpoint2(c *fiber.Ctx) error {
    // Your handler logic for endpoint2
    return c.SendString("Handler for test2")
}
