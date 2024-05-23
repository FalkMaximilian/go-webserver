package main

import (
	"go-webserver/config"
	"go-webserver/database"
	"go-webserver/logger"
	"go-webserver/middleware"
	"go-webserver/routes"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var testVariable string = "test"

func main() {

	var err error
	if err = config.SetupConfig(); err != nil {
		logger.Log.Fatal("Critical error: ", err)
	}
	logger.Log.Info("Reading config was successful!")

	// Setup fiber and cors
	app := fiber.New()
	app.Use(cors.New())
	app.Use(middleware.LoggingMiddleware())

	// Read port from env
	var port string = os.Getenv("PORT")
	if _, err = strconv.ParseUint(port, 10, 32); err != nil {
		logger.Log.Fatal("Critical error: environment variable 'PORT' must be set to a valid and unused port")
	}

	if err = database.ConnectDB(); err != nil {
		logger.Log.Fatal("Failed to connect to database: ", err)
	}

	// Routes without authentication

	routes.SetupAuthRoutes(app)
	routes.SetupAPIRoutes(app)

	// Start server
	logger.Log.Fatal(app.Listen(":" + port))
}
