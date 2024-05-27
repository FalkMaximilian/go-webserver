package handlers

import (
	"encoding/json"
	"fmt"
	"go-webserver/database"
	"go-webserver/logger"
	"go-webserver/models"
	"go-webserver/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserHandler(c *fiber.Ctx) error {

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		logger.Log.Info(err)
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		return fmt.Errorf("cannot parse json")
	}

	password, ok := data["password"].(string)
	if !ok {
		logger.Log.Info("password missing in request")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password is required"})
		return fmt.Errorf("password missing in request")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
		return fmt.Errorf("could not hash password")
	}

	data["password"] = string(hashedPassword)
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Log.Error(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
		return fmt.Errorf("could not marshal hashed password into json")
	}

	user := new(models.User)
	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		logger.Log.Error(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
		return fmt.Errorf("could not unmarshal json")
	}

	logger.Log.Println("Creating used in DB")
	result := database.DB.Create(user)
	if result.Error != nil {
		logger.Log.Error(result.Error)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
		return fmt.Errorf("error while creating user in db")
	}

	t, err := utils.GetJwtToken(user.ID)
	if err != nil {
		logger.Log.Error(err)
		c.SendStatus(fiber.StatusInternalServerError)
		return fmt.Errorf("error while generating jwt token")
	}

	logger.Log.Info("user registration successful")
	c.JSON(fiber.Map{"token": t})
	return fmt.Errorf("user registration successful")
}

func LoginHandler(c *fiber.Ctx) error {

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		return fmt.Errorf("could not parse json")
	}

	username, ok := data["username"].(string)
	if !ok {
		logger.Log.Warn("username misssing in request")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username is required"})
		return fmt.Errorf("username missing in request")
	}

	password, ok := data["password"].(string)
	if !ok {
		logger.Log.Warn("password misssing in request")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "passsword is required"})
		return fmt.Errorf("password missing in request")
	}

	user := new(models.User)
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		logger.Log.Info(err)
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "username not found"})
		return err
	}

	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "password incorrect"})
		return fmt.Errorf("incorrect password")
	}

	// Generate encoded token and send it as response.
	t, err := utils.GetJwtToken(user.ID)
	if err != nil {
		logger.Log.Error(err)
		c.SendStatus(fiber.StatusInternalServerError)
		return fmt.Errorf("error while generating jwt token")
	}

	c.JSON(fiber.Map{"token": t})
	return fmt.Errorf("login successful")
}

func DeleteUserHandler(c *fiber.Ctx) error {

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
		return fmt.Errorf("could not parse json")
	}

	username, ok := data["username"].(string)
	if !ok {
		logger.Log.Warn("username misssing in request")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username is required"})
		return fmt.Errorf("username missing in request")
	}

	password, ok := data["password"].(string)
	if !ok {
		logger.Log.Warn("password misssing in request")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passsword is required"})
		return fmt.Errorf("password missing in request")
	}

	user_id := c.Locals("user_id").(uint)

	user := new(models.User)
	if err := database.DB.Where("username = ? AND id = ?", username, user_id).First(&user).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete user"})
		return err
	}

	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil {
		logger.Log.Info("provided password incorrect")
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "provided password incorrect"})
		return fmt.Errorf("incorrect password")
	}

	// Load all sets for the user
	var sets []models.Set
	if err := database.DB.Where("user_id = ?", user_id).Find(&sets).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		return err
	}

	// Clear the associations for each set
	for _, set := range sets {
		if err := database.DB.Model(&set).Association("Cards").Clear(); err != nil {
			logger.Log.Warn(err)
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
			return err
		}
	}

	var card models.Card
	if err := database.DB.Delete(&card, "user_id = ?", user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		return err
	}

	var set models.Set
	if err := database.DB.Delete(&set, "user_id = ?", user_id).Error; err != nil {
		logger.Log.Warn(err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		return err
	}

	if err := database.DB.Model(&user).Association("Cards").Clear(); err != nil {
		logger.Log.Info("could not clear cards association")
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		return err
	}

	if err := database.DB.Model(&user).Association("Sets").Clear(); err != nil {
		logger.Log.Warn("could not clear sets association")
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		return err
	}

	if err := database.DB.Delete(&user, "id = ? AND username = ?", user_id, username).Error; err != nil {
		logger.Log.Info(err)
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "could not delete user"})
		return err
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{"success": "user deleted successfully"})
	return fmt.Errorf("user deleted successfully")
}
