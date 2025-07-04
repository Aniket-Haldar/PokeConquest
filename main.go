package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	initDatabase()
	SeedPokemon()

	router := mux.NewRouter()
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

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := ":8080"
	log.Printf("🚀 Server running at http://localhost%s", port)
	if err := http.ListenAndServe(port, corsHandler(router)); err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}
