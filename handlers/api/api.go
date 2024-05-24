package handlers

import (
	"fmt"
	"go-webserver/database"
	"go-webserver/logger"
	"go-webserver/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ProtectedTestHandler(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func ReadCards(c *fiber.Ctx) error {
	return nil
}

func CreateSet(c *fiber.Ctx) error {
	logger.Log.Info("Entered CreateSet function")

	// Get the username from Locals (set in auth middleware) and cast to uint
	user_id := c.Locals("user_id").(uint)
	logger.Log.WithField("user_id", user_id).Info("Extracted user ID from context")

	// Get all matched records
	var user models.User
	var err error
	if err = database.DB.First(&user, user_id).Error; err != nil {
		logger.Log.WithFields(logrus.Fields{
			"user_id": user_id,
			"error":   err,
		}).Error("Failed to find user")
		return fmt.Errorf("failed to find user: %v", err)
	}
	logger.Log.WithField("user_id", user_id).Info("User found in database")

	s := new(models.Set)
	if err = c.BodyParser(s); err != nil {
		return fmt.Errorf("failed to parse body into set: %v", err)
	}
	s.UserID = user_id

	// Save the set to the database
	if err = database.DB.Create(&s).Error; err != nil {
		return fmt.Errorf("failed to save set: %v", err)
	}
	logger.Log.WithFields(logrus.Fields{
		"user_id": user_id,
		"set":     s,
	}).Info("Set saved to database")

	// TODO: Name from parameter
	if err = database.DB.Model(&user).Association("Sets").Append(&s); err != nil {
		return fmt.Errorf("failed to append set to user: %v", err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func CreateCard(c *fiber.Ctx) error {
	return nil
}
