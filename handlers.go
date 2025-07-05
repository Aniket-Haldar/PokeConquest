package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
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
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // Earth radius in km
	latRad1 := lat1 * (math.Pi / 180)
	latRad2 := lat2 * (math.Pi / 180)
	deltaLat := (lat2 - lat1) * (math.Pi / 180)
	deltaLng := (lng2 - lng1) * (math.Pi / 180)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latRad1)*math.Cos(latRad2)*math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func getNearbyTrainers(w http.ResponseWriter, r *http.Request) {
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

	// Find all users within 2km (you can adjust the radius)
	var trainers []User
	if err := DB.Find(&trainers).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	nearby := []User{}
	for _, t := range trainers {
		if t.LastLocation.Lat != 0 && t.LastLocation.Lng != 0 {
			distance := haversineDistance(lat, lng, t.LastLocation.Lat, t.LastLocation.Lng)
			if distance <= 2.0 {
				nearby = append(nearby, t)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nearby)
}

// Helper: Haversine distance between two lat/lng
// Load API Key from ENV

const geminiURL = "https://generativelanguage.googleapis.com/v1/models/gemini-2.0-flash:generateContent"

// Request & Response Structs
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Send prompt to Gemini API and get response
func getGeminiResponse(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Gemini API key not set in environment")
	}

	// Send API key as query param
	url := fmt.Sprintf("%s?key=%s", geminiURL, apiKey)

	body := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// REMOVE this line: req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Gemini API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return "", fmt.Errorf("Gemini API error: %v", errResp)
	}

	var result GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
func trimText(text string, maxWords int) string {
	words := bytes.Fields([]byte(text))
	if len(words) <= maxWords {
		return text
	}
	return string(bytes.Join(words[:maxWords], []byte(" "))) + "..."
}

func aiTipHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Context string `json:"context"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	prompt := "Give a short and helpful tip (max 2 lines) about: " + req.Context
	tip, err := getGeminiResponse(prompt)
	if err != nil {
		http.Error(w, fmt.Sprintf("AI error: %v", err), http.StatusInternalServerError)
		return
	}

	tip = trimText(tip, 60)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"tip": tip,
	})
}

func aiStrategyHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pokemon string `json:"pokemon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	prompt := fmt.Sprintf("Give a battle strategy for using %s in Pokémon battles.", req.Pokemon)
	strategy, err := getGeminiResponse(prompt)
	if err != nil {
		http.Error(w, fmt.Sprintf("AI error: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"strategy": strategy})
}
