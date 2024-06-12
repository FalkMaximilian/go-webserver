package utils

import (
	"go-webserver/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Create a new jwt token that encodes the user_id
func GetJwtToken(user_id uint, username string) (string, error) {

	claims := jwt.MapClaims{
		"user_id":  user_id,
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.GetJWTSecret())
}
