package middleware

import (
	"fmt"
	"go-webserver/config"
	"go-webserver/logger"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

func JWTProtected() fiber.Handler {

	return func(c *fiber.Ctx) error {

		// Get the token from the request header
		tokenString := c.Get("Authorization")
		logger.Log.Debug(tokenString)

		if tokenString == "" {
			logger.Log.Warn("missing jwt token")
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing jwt token"})
			return fmt.Errorf("missing jwt token")
		}

		// Remove 'Bearer: ' prefix if present
		if strings.HasPrefix(tokenString, "Bearer ") {
			logger.Log.Debug("removing 'Bearer ' prefix from jwt token")
			tokenString = tokenString[len("Bearer "):]
		}

		// Parse the jwt token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.GetJWTSecret(), nil
		})

		if err != nil {
			logger.Log.WithField("error", err).Warn("invalid jwt token")
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid jwt token"})
			return err
		}

		// Extract user information from the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user := claims["user_id"].(float64) // Get user_id as f64
			c.Locals("user_id", uint(user))     // Store the user ID in the context locals as uint
		} else {
			logger.Log.Warn("invalid jwt token")
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid jwt token"})
			return fmt.Errorf("invalid jwt token")
		}

		return c.Next()
	}
}

func RedirectIfAuthenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			logger.Log.Debug("no jwt token - continue")
			return c.Next()
		}

		// Remove 'Bearer: ' prefix if present
		if strings.HasPrefix(tokenString, "Bearer ") {
			logger.Log.Debug("removing 'Bearer ' prefix from jwt token")
			tokenString = tokenString[len("Bearer "):]
		}

		// Parse the jwt token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.GetJWTSecret(), nil
		})

		if err == nil && token.Valid {
			c.Redirect("/")
			return fmt.Errorf("redirecting to home")
		}

		return c.Next()
	}
}

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		if sc := c.Response().StatusCode(); sc >= 500 {
			logger.Log.WithFields(logrus.Fields{
				"status":     c.Response().StatusCode(),
				"method":     c.Method(),
				"path":       c.Path(),
				"latency":    stop.Sub(start).String(),
				"client_ip":  c.IP(),
				"user_agent": c.Get("User-Agent"),
			}).Error(err)
		} else if sc >= 400 {
			logger.Log.WithFields(logrus.Fields{
				"status":     c.Response().StatusCode(),
				"method":     c.Method(),
				"path":       c.Path(),
				"latency":    stop.Sub(start).String(),
				"client_ip":  c.IP(),
				"user_agent": c.Get("User-Agent"),
			}).Warn(err)
		} else {
			logger.Log.WithFields(logrus.Fields{
				"status":     c.Response().StatusCode(),
				"method":     c.Method(),
				"path":       c.Path(),
				"latency":    stop.Sub(start).String(),
				"client_ip":  c.IP(),
				"user_agent": c.Get("User-Agent"),
			}).Info(err)
		}

		return nil
	}
}
