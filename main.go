package main

import (
	"encoding/json"
	"go-webserver/database"
	"go-webserver/model"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Fiber instance
	app := fiber.New()

	var port string = os.Getenv("PORT")

	database.ConnectDB()

	// Routes
	app.Get("/", hello)
	app.Post("/register", registerUser)

	// Start server
	log.Fatal(app.Listen(":" + port))
}

func registerUser(c *fiber.Ctx) error {

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

	user := new(model.User)
	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to has password"})
	}

	log.Println("Creating used in DB")
	database.DB.Create(user)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": "User was created"})
}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
