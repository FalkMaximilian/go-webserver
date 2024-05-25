package handlers

import (
	"fmt"
	"go-webserver/database"
	"go-webserver/logger"
	"go-webserver/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func CreateSetHandler(c *fiber.Ctx) error {
	logger.Log.Debug("entering 'CreateSet' handler")

	// Get the username from Locals (set in auth middleware) and cast to uint
	user_id := c.Locals("user_id").(uint)
	//logger.Log.WithField("user_id", user_id).Info("Extracted user ID from context")

	// Get all matched records
	var user models.User
	var err error
	if err = database.DB.First(&user, user_id).Error; err != nil {
		logger.Log.WithFields(logrus.Fields{
			"user_id": user_id,
			"error":   err,
		}).Error("unable to find user")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create set"})
		return fmt.Errorf("could not find user in db")
	}
	logger.Log.WithField("user_id", user_id).Debug("user found in database")

	s := new(models.Set)
	if err = c.BodyParser(s); err != nil {
		logger.Log.WithField("error", err).Error("failed to parse request body")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create set"})
		return fmt.Errorf("could not parse request body into set")
	}
	s.UserID = user_id

	// Save the set to the database
	if err = database.DB.Create(&s).Error; err != nil {
		logger.Log.WithField("error", err).Error("failed to create set")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create set"})
		return fmt.Errorf("could not create set")
	}

	if err = database.DB.Model(&user).Association("Sets").Append(&s); err != nil {
		// TODO: What happens if this failes but the previous one does not fail???
		logger.Log.WithField("error", err).Error("could not append set to user")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create set"})
		return fmt.Errorf("could not append new set to user association")
	}

	c.Status(fiber.StatusCreated).JSON(fiber.Map{"set": s})
	return fmt.Errorf("set created successfully")
}

func DeleteSetHandler(c *fiber.Ctx) error {
	set_id := c.Params("set_id")
	user_id := c.Locals("user_id").(uint)

	var set models.Set
	if err := database.DB.Find(&set, "id = ? AND user_id = ?", set_id, user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "could not delete set"})
		return err
	}

	if err := database.DB.Model(&set).Association("Cards").Clear(); err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete set"})
		return err
	}

	if err := database.DB.Delete(&set).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete set"})
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "set deleted successfully"})
	return fmt.Errorf("set deleted successfully")
}

func AddCardToSetHandler(c *fiber.Ctx) error {

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		logger.Log.Info(err)
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		return fmt.Errorf("cannot parse json")
	}

	set_id := c.Params("set_id")
	user_id := c.Locals("user_id").(uint)
	card_id := data["card_id"]

	// Set has to exist and be owned by user
	var set models.Set
	if err := database.DB.Find(&set, "id = ? AND user_id = ?", set_id, user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "set could not be found"})
		return err
	}

	// Card has to exist and be owned by user
	var card models.Card
	if err := database.DB.Find(&card, "id = ? AND user_id = ?", card_id, user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "card could not be found"})
		return err
	}

	if err := database.DB.Model(&set).Association("Cards").Append(&card); err != nil {
		logger.Log.Error(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "card could not be added to set"})
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"set_id": set_id, "card_id": card_id})
	return fmt.Errorf("card successfully added to set")
}
