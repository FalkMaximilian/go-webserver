package handlers

import (
	"encoding/json"
	"fmt"
	"go-webserver/database"
	"go-webserver/models"
	"go-webserver/token"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserHandler(c *fiber.Ctx) error {
	// Check if user is already signed in
	tokenString := c.Get("Authorization")
	if len(tokenString) > 7 {
		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			log.Println(token.Method.Alg())
			if alg := token.Method.Alg(); alg != "HS256" {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err == nil && token.Valid {
			log.Println("Redirect to /")
			return c.Redirect("/")
		}
	}

	log.Println("Parsing input...")
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	password, ok := data["password"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password is required"})
	}

	log.Println("Bcrypt password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to has password"})
	}

	data["password"] = string(hashedPassword)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to has password"})
	}

	user := new(models.User)
	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to has password"})
	}

	log.Println("Creating used in DB")
	result := database.DB.Create(user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Generate encoded token and send it as response.
	t, err := token.GetJwtToken(user.Username)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func LoginHandler(c *fiber.Ctx) error {

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	username, ok := data["username"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username is required"})
	}

	password, ok := data["password"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passsword is required"})
	}

	user := new(models.User)
	database.DB.Where("username = ?", username).First(&user)
	log.Println(user)

	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Password is wrong!"})
	}

	// Generate encoded token and send it as response.
	t, err := token.GetJwtToken(user.Username)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
