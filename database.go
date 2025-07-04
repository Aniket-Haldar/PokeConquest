package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDatabase() {

	dsn := "host=localhost user=pokeuser password=pokepass dbname=pokequest port=5432 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
		os.Exit(1)
	}

	err = DB.AutoMigrate(&User{}, &Pokemon{}, &Challenge{})
	if err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	log.Println("✅ Database connected and migrated successfully")
}
