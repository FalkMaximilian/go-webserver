package main

import (
	"log"
	"os"

	"go-webserver/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Fiber instance
	app := fiber.New()

	var port string = os.Getenv("PORT")

	database.ConnectDB()

	// Routes
	app.Get("/", hello)

	// Start server
	log.Fatal(app.Listen(":" + port))
}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
