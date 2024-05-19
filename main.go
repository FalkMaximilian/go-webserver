package main

import (
	"encoding/json"
	"fmt"
	"go-webserver/api"
	"go-webserver/auth"
	"go-webserver/database"
	"go-webserver/model"
	"log"
	"os"
	"strconv"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/joho/godotenv"
)

func main() {

	// Read environment vars for local development
	err := godotenv.Load()
	if err != nil {
		log.Printf("Could not load .env file: %v", err)
	}

	// Setup fiber and cors
	app := fiber.New()
	app.Use(cors.New())

	// Read port from env
	var port string = os.Getenv("PORT")
	if _, err = strconv.ParseUint(port, 10, 32); err != nil {
		log.Fatal("Critical error: environment variable 'PORT' must be set to a valid and unused port")
	}

	if err = database.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = auth.SetupJWT(); err != nil {
		log.Fatalf("Failed to set up JWT: %v", err)
	}

	// Routes without authentication
	app.Get("/", hello)
	app.Post("/register", registerUser)
	app.Post("/login", login)

	app.Use(jwtware.New(auth.GetConfig()))

	api_handlers := app.Group("/api")
	api.RegisterHandlers(api_handlers)

	// Start server
	log.Fatal(app.Listen(":" + port))
}

func registerUser(c *fiber.Ctx) error {
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

	user := new(model.User)
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
	t, err := auth.GetJWT(user.Username)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
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

	// Generate encoded token and send it as response.
	t, err := auth.GetJWT(user.Username)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

// Handler
func hello(c *fiber.Ctx) error {
	log.Println("Hello endpoint called!")
	return c.SendString("Hello, World 👋!")
}
