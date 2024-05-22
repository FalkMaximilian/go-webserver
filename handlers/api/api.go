package handlers

import (
	"fmt"
	"go-webserver/database"
	"go-webserver/models"

	"github.com/gofiber/fiber/v2"
)

func ProtectedTestHandler(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func ReadCards(c *fiber.Ctx) error {
	return nil
}

func CreateSet(c *fiber.Ctx) error {

	// Get the username from Locals (set in auth middleware)
	user_id := c.Locals("user_id")

	// Get all matched records
	var user models.User
	if err := database.DB.First(&user, user_id).Error; err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	// TODO: Name from parameter
	if err := database.DB.Model(&user).Association("Sets").Append(&models.Set{Name: "Test123"}); err != nil {
		return fmt.Errorf("failed to append set to user: %v", err)
	}

	return c.SendStatus(fiber.StatusCreated)
}
