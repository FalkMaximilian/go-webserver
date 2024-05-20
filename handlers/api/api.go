package handlers

import "github.com/gofiber/fiber/v2"

func ProtectedTestHandler(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
