package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Fiber instance
	app := fiber.New()

  var port string = os.Getenv("PORT")

	// Routes
	app.Get("/", hello)

	// Start server
	log.Fatal(app.Listen(":" + port))
}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
