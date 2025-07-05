package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file not found. Using defaults.")
	}

	// Initialize DB and seed data
	initDatabase()
	SeedPokemon()

	// Create router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", createUser).Methods("POST")
	api.HandleFunc("/users/{id}", getUser).Methods("GET")
	api.HandleFunc("/users/{id}/location", updateUserLocation).Methods("PUT")
	api.HandleFunc("/users/{id}/gamestate", getGameState).Methods("GET")

	api.HandleFunc("/challenges", generateChallenges).Methods("GET")
	api.HandleFunc("/challenges/{id}/complete", completeChallenge).Methods("POST")

	api.HandleFunc("/pokemon/nearby", getNearbyPokemon).Methods("GET")
	api.HandleFunc("/locations/pois", getNearbyPOIs).Methods("GET")

	api.HandleFunc("/health", healthCheck).Methods("GET")
	api.HandleFunc("/catch", catchPokemon).Methods("POST")
	api.HandleFunc("/trainers/nearby", getNearbyTrainers).Methods("GET")
	api.HandleFunc("/ai/tip", aiTipHandler).Methods("POST")
	api.HandleFunc("/ai/strategy", aiStrategyHandler).Methods("POST")

	// Serve frontend static files
	frontendDir := "./frontend"
	fs := http.FileServer(http.Dir(frontendDir))
	router.PathPrefix("/").Handler(fs)

	// Enable CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Get port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, corsHandler(router)); err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}
