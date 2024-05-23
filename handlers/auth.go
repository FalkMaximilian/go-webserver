package handlers

import (
	"encoding/json"
	"fmt"
	"go-webserver/database"
	"go-webserver/logger"
	"go-webserver/models"
	"go-webserver/token"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserHandler(c *fiber.Ctx) error {

	// TODO: Add new middleware that handles redirection for these cases
	// Check if user is already signed in
	tokenString := c.Get("Authorization")
	if len(tokenString) > 7 {
		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			if alg := token.Method.Alg(); alg != "HS256" {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err == nil && token.Valid {
			logger.Log.Info("Redirecting to /")
			return c.Redirect("/")
		}
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		logger.Log.Debug("RegisterUserHandler: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	password, ok := data["password"].(string)
	if !ok {
		logger.Log.Debug("RegisterUserHandler: Password missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password is required"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("RegisterUserhandler: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	data["password"] = string(hashedPassword)
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Log.Error("RegisterUserHandler: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := new(models.User)
	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		logger.Log.Error("RegisterUserHandler: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	logger.Log.Println("Creating used in DB")
	result := database.DB.Create(user)
	if result.Error != nil {
		logger.Log.Error("RegisterUserHandler: ", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Generate encoded token and send it as response.
	t, err := token.GetJwtToken(user.ID)
	if err != nil {
		logger.Log.Error("RegisterUserHandler: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	logger.Log.Info("RegisterUserHandler: User registration successful. Issuing auth token (jwt)")
	return c.JSON(fiber.Map{"token": t})
}

func LoginHandler(c *fiber.Ctx) error {

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		logger.Log.Debug("logger.Loginhandler: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	username, ok := data["username"].(string)
	if !ok {
		logger.Log.Debug("logger.LoginHandler: username misssing in logger.Login request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username is required"})
	}

	password, ok := data["password"].(string)
	if !ok {
		logger.Log.Debug("logger.LoginHandler: password misssing in logger.Login request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passsword is required"})
	}

	user := new(models.User)
	database.DB.Where("username = ?", username).First(&user)

	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil {
		logger.Log.Info("logger.LoginHandler: Failed logger.Login attempt. Wrong password.")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Password is wrong!"})
	}

	// Generate encoded token and send it as response.
	t, err := token.GetJwtToken(user.ID)
	if err != nil {
		logger.Log.Error("logger.Loginhandler: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	logger.Log.Info("logger.LoginHandler: logger.Login successful. Issuing auth token (jwt)")
	return c.JSON(fiber.Map{"token": t})
}
