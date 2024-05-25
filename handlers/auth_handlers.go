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
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username is required"})
		return fmt.Errorf("username missing in request")
	}

	password, ok := data["password"].(string)
	if !ok {
		logger.Log.Warn("password misssing in request")
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passsword is required"})
		return fmt.Errorf("password missing in request")
	}

	user := new(models.User)
	database.DB.Where("username = ?", username).First(&user)

	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Password is wrong!"})
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
