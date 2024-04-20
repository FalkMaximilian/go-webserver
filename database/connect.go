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

func ConnectDB() {

	// Get database port and test if correct
	db_port := os.Getenv("PGPORT")
	_, err := strconv.ParseUint(db_port, 10, 32)

	if err != nil {
		log.Fatal("PGPORT invalid! Abort...")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"),
		db_port,
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGNAME"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not establish database connection! Abort...")
	}

	log.Println("Established connection to database.")
	DB.AutoMigrate(&model.User{})
	log.Println("Database migrated")
}
