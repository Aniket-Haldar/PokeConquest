package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDatabase() {
	// Load DB settings from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Fallback defaults for local dev
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "pokeuser"
	}
	if dbPass == "" {
		dbPass = "pokepass"
	}
	if dbName == "" {
		dbName = "pokequest"
	}

	// Connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

	// Connect to DB
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database: ", err)
		os.Exit(1)
	}

	// Auto-migrate schema
	err = DB.AutoMigrate(&User{}, &Pokemon{}, &Challenge{})
	if err != nil {
		log.Fatal("❌ Failed to run migrations: ", err)
	}

	log.Println("✅ Database connected and migrated successfully")
}
