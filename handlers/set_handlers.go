package handlers

import (
	"fmt"
	"go-webserver/database"
	"go-webserver/logger"
	"go-webserver/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Get all sets for the user requesting it
func GetSetsHandler(c *fiber.Ctx) error {
	logger.Log.Debug("entering 'GetSets' handler")

	user_id := c.Locals("user_id").(uint)

	var sets []models.Set
	if err := database.DB.Find(&sets, "user_id = ?", user_id).Error; err != nil {
		logger.Log.Info(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve sets"})
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"sets": sets})
	return fmt.Errorf("sets retreived successfully")
}

// Get set with id. Has to belong to user requesting it
func GetSetHandler(c *fiber.Ctx) error {
	logger.Log.Debug("entering 'GetSet' handler")

	user_id := c.Locals("user_id").(uint)
	set_id := c.Params("set_id")

	// Get set with id. Has to be owned by user requesting it
	var set models.Set
	if err := database.DB.First(&set, "id = ? AND user_id = ?", set_id, user_id).Error; err != nil {
		logger.Log.Info(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get set"})
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"set": set})
	return fmt.Errorf("set retreived successfully")
}

// Create a new set
func CreateSetHandler(c *fiber.Ctx) error {
	logger.Log.Debug("entering 'CreateSet' handler")

	// Get the username
	user_id := c.Locals("user_id").(uint)

	// Get user for user_id
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

	// Parse request body into set
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

// Remove all cards from a set and then delete the set
func DeleteSetHandler(c *fiber.Ctx) error {
	set_id := c.Params("set_id")
	user_id := c.Locals("user_id").(uint)

	// Get set by id. Has to belong to the user requesting it
	var set models.Set
	if err := database.DB.Find(&set, "id = ? AND user_id = ?", set_id, user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "could not delete set"})
		return err
	}

	// Clear all associations to Cards for that set (remove entries from set_cards)
	if err := database.DB.Model(&set).Association("Cards").Clear(); err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete set"})
		return err
	}

	// Finally delete the set
	if err := database.DB.Delete(&set).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete set"})
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "set deleted successfully"})
	return fmt.Errorf("set deleted successfully")
}

// Used to add a card to a set. Both have to be owned by the user requesting said action
func AddCardToSetHandler(c *fiber.Ctx) error {

	// Parse json body into data
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

	// Add card to set by setting association in set_cards
	if err := database.DB.Model(&set).Association("Cards").Append(&card); err != nil {
		logger.Log.Error(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "card could not be added to set"})
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"set_id": set_id, "card_id": card_id})
	return fmt.Errorf("card successfully added to set")
}

// Not yet implemented
func RemoveCardFromSetHandler(c *fiber.Ctx) error {
	return fmt.Errorf("card successfully removed from set")
}
