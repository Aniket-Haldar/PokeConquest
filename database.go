package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDatabase() {
	// Check for DATABASE_URL (Render production)
	dbURL := os.Getenv("DATABASE_URL")

	// If DATABASE_URL is not set, build from individual vars (local dev)
	if dbURL == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		// Fallback defaults for local development
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

		dbURL = "host=" + dbHost +
			" user=" + dbUser +
			" password=" + dbPass +
			" dbname=" + dbName +
			" port=" + dbPort +
			" sslmode=disable"
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Assign global DB
	DB = db

	// Auto-migrate schema
	err = DB.AutoMigrate(&User{}, &Pokemon{}, &Challenge{})
	if err != nil {
		log.Fatal("❌ Failed to run migrations: ", err)
	}

	log.Println("✅ Database connected and migrated successfully")
}
