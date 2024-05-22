package token

import (
	"go-webserver/config"
	"go-webserver/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetJwtToken(user_id uint) (string, error) {
	logger.Debug("Issuing new JWT token")

	claims := jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.GetJWTSecret())
}

/*
func GetConfig() jwtware.Config {
	return jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.GetJWTSecret())},
	}
}
*/
