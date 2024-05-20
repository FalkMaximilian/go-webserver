package token

import (
	"go-webserver/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetJwtToken(name string) (string, error) {
	claims := jwt.MapClaims{
		"name":  name,
		"admin": false,
		"exp":   time.Now().Add(time.Hour * 24 * time.Duration(config.GetJWTExpires())).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.GetJWTSecret()))
}

/*
func GetConfig() jwtware.Config {
	return jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.GetJWTSecret())},
	}
}
*/
