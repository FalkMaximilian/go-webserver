package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"go-webserver/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() error {

	var (
		pghost     string
		pgport     string
		pguser     string
		pgpassword string
		pgname     string
	)

	if pghost = os.Getenv("PGHOST"); pghost == "" {
		return fmt.Errorf("environment variable 'PGHOST' missing")
	}

	pgport = os.Getenv("PGPORT")
	if _, err := strconv.ParseUint(pgport, 10, 32); err != nil {
		return fmt.Errorf("environment variable 'PGPORT' missing or invalid")
	}

	if pguser = os.Getenv("PGUSER"); pguser == "" {
		return fmt.Errorf("environment variable 'PGUSER' missing")
	}

	if pgpassword = os.Getenv("PGPASSWORD"); pgpassword == "" {
		return fmt.Errorf("environment variable 'PGPASSWORD' missing")
	}

	if pgname = os.Getenv("PGNAME"); pgname == "" {
		return fmt.Errorf("environment variable 'PGNAME' missing")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pghost,
		pgport,
		pguser,
		pgpassword,
		pgname,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	DB.AutoMigrate(&model.User{})
	log.Printf("Database setup successful")
	return nil
}
