package middleware

import (
	"fmt"
	"go-webserver/config"
	"go-webserver/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

func JWTProtected() fiber.Handler {

	return func(c *fiber.Ctx) error {

		// Get the token from the request header
		tokenString := c.Get("Authorization")
		if len(tokenString) <= 7 {
			// logger.Warn("Missing or invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid JWT token",
			})
		}

		// Remove 'bearer' from token
		tokenString = tokenString[len("Bearer "):]

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.GetJWTSecret(), nil
		})

		if err != nil {
			// logger.Warn("Invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid JWT token",
			})
		}

		// Extract user information from the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user := claims["user_id"].(float64) // Get user_id as f64
			c.Locals("user_id", uint(user))     // Store the user ID in the context locals as uint
			// logger.Debug("Added 'user_id' to Locals")
		} else {
			// logger.Warn("Invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid JWT token",
			})
		}
		return c.Next()
	}
}

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		logger.Log.WithFields(logrus.Fields{
			"status":     c.Response().StatusCode(),
			"method":     c.Method(),
			"path":       c.Path(),
			"latency":    stop.Sub(start).String(),
			"client_ip":  c.IP(),
			"user_agent": c.Get("User-Agent"),
		}).Info("request completed")

		return err
	}
}
