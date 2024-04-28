package main

import (
	"encoding/json"
	"fmt"
	"go-webserver/database"
	"go-webserver/model"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Fiber instance
	app := fiber.New()
	app.Use(cors.New())

	var port string = os.Getenv("PORT")

	database.ConnectDB()

	// Routes
	app.Get("/", hello)
	app.Post("/register", registerUser)
	app.Post("/login", login)

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
	result := database.DB.Create(user)
	log.Println(result.Error)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": "User was created"})
}

func login(c *fiber.Ctx) error {

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

	user := new(model.User)
	database.DB.Where("username = ?", username).First(&user)
	log.Println(user)

	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if result != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Password is wrong!"})
	}

	return c.SendString(fmt.Sprintf("Das Passwort für User '%s' war korrekt", username))
}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
