package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := User{
		Username: req.Username,
		Level:    1,
		XP:       0,
	}
	if err := DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var user User
	if err := DB.Preload("PokemonCaught").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(user)
}

func updateUserLocation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req Location
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.LastLocation = req
	if err := DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update location", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func getGameState(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var user User
	if err := DB.Preload("PokemonCaught").First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var activeChallenges []Challenge
	DB.Where("completed = ?", false).Find(&activeChallenges)

	resp := map[string]interface{}{
		"user":              user,
		"active_challenges": activeChallenges,
	}
	json.NewEncoder(w).Encode(resp)
}

func generateChallenges(w http.ResponseWriter, r *http.Request) {
	var challenges []Challenge
	for i := 0; i < 3; i++ {
		pokemon, _ := selectRandomPokemon()
		challenge := Challenge{
			Title:       fmt.Sprintf("Catch %s", pokemon.Name),
			Description: "Find and catch this Pokémon!",
			Pokemon:     pokemon,
			Completed:   false,
			Location:    generateRandomLocation(22.5726, 88.3639, 2), // random near Kolkata
		}
		DB.Create(&challenge)
		challenges = append(challenges, challenge)
	}
	json.NewEncoder(w).Encode(challenges)
}

func completeChallenge(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var challenge Challenge
	if err := DB.First(&challenge, id).Error; err != nil {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	challenge.Completed = true
	DB.Save(&challenge)

	var user User
	DB.First(&user, req.UserID)
	user.XP += 100
	user.PokemonCaught = append(user.PokemonCaught, challenge.Pokemon)
	DB.Save(&user)

	resp := map[string]interface{}{
		"success":   true,
		"user":      user,
		"challenge": challenge,
	}
	json.NewEncoder(w).Encode(resp)
}

func getNearbyPokemon(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	var spawns []map[string]interface{}

	for i := 0; i < 5; i++ {
		pokemon, err := selectRandomPokemon()
		if err != nil {
			http.Error(w, "Could not fetch Pokémon", http.StatusInternalServerError)
			return
		}

		location := generateRandomLocation(lat, lng, 1.0)
		spawn := map[string]interface{}{
			"pokemon":      pokemon,
			"location":     location,
			"despawn_time": time.Now().Add(5 * time.Minute).Unix(),
		}

		spawns = append(spawns, spawn)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spawns)
}

func getNearbyPOIs(w http.ResponseWriter, r *http.Request) {
	pois := []map[string]interface{}{
		{"name": "Park", "type": "park", "lat": 22.5726, "lng": 88.3639},
		{"name": "Museum", "type": "museum", "lat": 22.5730, "lng": 88.3645},
	}
	json.NewEncoder(w).Encode(pois)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func catchPokemon(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	pokemonID, err := strconv.Atoi(r.URL.Query().Get("pokemon_id"))
	if err != nil {
		http.Error(w, "Invalid Pokémon ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := DB.Preload("PokemonCaught").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	var pokemon Pokemon
	if err := DB.First(&pokemon, pokemonID).Error; err != nil {

		pokemon, err = fetchPokemonFromAPI(pokemonID)
		if err != nil {
			http.Error(w, "Failed to fetch Pokémon data", http.StatusInternalServerError)
			return
		}
		if err := DB.Create(&pokemon).Error; err != nil {
			http.Error(w, "Failed to save Pokémon", http.StatusInternalServerError)
			return
		}
	}

	if err := DB.Model(&user).Association("PokemonCaught").Append(&pokemon); err != nil {
		http.Error(w, "Failed to update Pokédex", http.StatusInternalServerError)
		return
	}

	if user.XP >= 100 {
		user.Level++
		user.XP = 0
	}

	if err := DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to save user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pokémon caught successfully!",
		"user":    user,
	})
}
