package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

var expires_in uint64
var jwt_secret string

func SetupJWT() error {
	if jwt_secret = os.Getenv("JWT_SECRET"); len(jwt_secret) <= 32 {
		return fmt.Errorf("environment variable 'JWT_SECRET' must be set and at least 32 characters long")
	}

	var err error
	expires_in, err = strconv.ParseUint(os.Getenv("JWT_EXPIRES_IN"), 10, 32)
	if err != nil {
		return fmt.Errorf("environment variable 'JWT_EXPIRES_IN' must be set to an unsigned int")
	}

	return nil
}

func GetJWT(name string) (string, error) {
	claims := jwt.MapClaims{
		"name":  name,
		"admin": false,
		"exp":   time.Now().Add(time.Hour * 24 * time.Duration(expires_in)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwt_secret))
}

func GetConfig() jwtware.Config {
	return jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwt_secret)},
	}
}
