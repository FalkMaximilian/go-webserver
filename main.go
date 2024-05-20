package main

import (
	"go-webserver/config"
	"go-webserver/database"
	"go-webserver/routes"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	var err error
	if err = config.SetupConfig(); err != nil {
		log.Fatalf("Critical error: %v", err)
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

	// Routes without authentication

	routes.SetupAuthRoutes(app)
	routes.SetupAPIRoutes(app)

	// Start server
	log.Fatal(app.Listen(":" + port))
}
