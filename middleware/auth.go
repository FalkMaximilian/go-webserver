package middleware

import (
	"fmt"
	"go-webserver/config"
	"go-webserver/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JWTProtected() fiber.Handler {

	return func(c *fiber.Ctx) error {

		// Get the token from the request header
		tokenString := c.Get("Authorization")
		if len(tokenString) <= 7 {
			logger.Warn("Missing or invalid JWT token")
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
			logger.Warn("Invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid JWT token",
			})
		}

		// Extract user information from the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user := claims["user_id"] // Adjust the key according to your JWT payload structure
			c.Locals("user_id", user) // Store the user ID in the context locals
			logger.Debug("Added 'user_id' to Locals")
		} else {
			logger.Warn("Invalid JWT token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid JWT token",
			})
		}
		return c.Next()
	}
}
