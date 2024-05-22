package main

import (
	"go-webserver/config"
	"go-webserver/database"
	"go-webserver/logger"
	"go-webserver/routes"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	var err error
	if err = config.SetupConfig(); err != nil {
		logger.Fatal("Critical error: ", err)
	}
	logger.Info("Reading config was successful!")

	// Setup fiber and cors
	app := fiber.New()
	app.Use(cors.New())

	// Read port from env
	var port string = os.Getenv("PORT")
	if _, err = strconv.ParseUint(port, 10, 32); err != nil {
		logger.Fatal("Critical error: environment variable 'PORT' must be set to a valid and unused port")
	}

	if err = database.ConnectDB(); err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}

	// Routes without authentication

	routes.SetupAuthRoutes(app)
	routes.SetupAPIRoutes(app)

	// Start server
	logger.Fatal(app.Listen(":" + port))
}
