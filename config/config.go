package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

type JwtConfig struct {
	secret     string
	expires_in uint64
}

var jwt_config JwtConfig

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

	// Read the JWT config
	if jwt_config.secret = os.Getenv("JWT_SECRET"); len(jwt_config.secret) <= 32 {
		return fmt.Errorf("environment variabel 'JWT_SECRET' must be at least 32 chars")
	}

	if jwt_config.expires_in, err = strconv.ParseUint(os.Getenv("JWT_EXPIRES_IN"), 10, 32); err != nil {
		return fmt.Errorf("environment variable 'JWT_EXPIRES_IN' must be set to an unsigned int")
	}

	return nil
}

func GetJWTSecret() string {
	return jwt_config.secret
}

func GetJWTExpires() uint64 {
	return jwt_config.expires_in
}
