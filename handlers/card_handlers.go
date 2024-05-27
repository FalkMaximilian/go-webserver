package handlers

import (
	"fmt"
	"go-webserver/database"
	"go-webserver/logger"
	"go-webserver/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetCardHandler(c *fiber.Ctx) error {
	logger.Log.Debug("entering 'GetCard' handler")

	user_id := c.Locals("user_id").(uint)
	card_id := c.Params("card_id")

	var card models.Card
	if err := database.DB.First(&card, "id = ? AND user_id = ?", card_id, user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get card"})
		return fmt.Errorf("failed to get card")
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"card": card})
	return fmt.Errorf("successfully received card")
}

func GetCardsHandler(c *fiber.Ctx) error {
	logger.Log.Debug("entering 'GetCard' handler")

	user_id := c.Locals("user_id").(uint)

	var cards []models.Card
	if err := database.DB.Find(&cards, "user_id = ?", user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get cards"})
		return fmt.Errorf("failed to get cards")
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"cards": cards})
	return fmt.Errorf("successfully received card")
}

func CreateCardHandler(c *fiber.Ctx) error {

	user_id := c.Locals("user_id").(uint)

	var user models.User
	var err error
	if err = database.DB.First(&user, user_id).Error; err != nil {
		logger.Log.WithFields(logrus.Fields{
			"user_id": user_id,
			"error":   err,
		}).Error("failed to find user")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create card"})
		return fmt.Errorf("could not find user in db")
	}

	card := new(models.Card)
	if err = c.BodyParser(card); err != nil {
		logger.Log.WithField("error", err).Error("failed to parse request body")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create card"})
		return fmt.Errorf("could not parse request body into card")
	}
	card.UserID = user_id

	if err = database.DB.Create(&card).Error; err != nil {
		logger.Log.WithField("error", err).Error("failed to create card")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create card"})
		return fmt.Errorf("could not create card")
	}

	if err = database.DB.Model(&user).Association("Cards").Append(&card); err != nil {
		logger.Log.WithField("error", err).Error("could not append card to user")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to create card"})
		return fmt.Errorf("could not append new card to user association")
	}

	c.Status(fiber.StatusCreated).JSON(fiber.Map{"card": card})
	return fmt.Errorf("card created successfully")
}

func DeleteCardHandler(c *fiber.Ctx) error {
	card_id := c.Params("card_id")
	user_id := c.Locals("user_id").(uint)

	var card models.Card
	if err := database.DB.First(&card, "id = ? AND user_id = ?", card_id, user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "card could not be deleted"})
		return err
	}

	if err := database.DB.Model(&card).Association("Sets").Clear(); err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "card could not be deleted"})
		return err
	}

	if err := database.DB.Delete(&card).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "card could not be deleted"})
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "card deleted successfully"})
	return fmt.Errorf("card deleted successfully")
}
