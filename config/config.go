package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
type PgConfig struct {
	pgname     string
	pghost     string
	pguser     string
	pgport     string
	pgpassword string
}

type FiberConfig struct {
	port string
}
*/

var jwt_secret_key []byte

/*
var PgConfig PgConfig
var FiberConfig FiberConfig
var JwtConfig JwtConfig
*/

func SetupConfig() error {

	// Read environment vars for local development
	var err error = godotenv.Load()
	if err != nil {
		log.Printf("Warning: %v", err)
	}

	// Read the JWT secret
	temp_sec := os.Getenv("JWT_SECRET")
	if len(temp_sec) < 32 {
		return fmt.Errorf("environment variabel 'JWT_SECRET' must be at least 32 chars")
	}
	jwt_secret_key = []byte(temp_sec)

	return nil
}

func GetJWTSecret() []byte {
	return jwt_secret_key
}
