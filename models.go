package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
)

const pokeapiURL = "https://pokeapi.co/api/v2/pokemon/"

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Pokemon struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Sprite string `json:"sprite"`
}

type User struct {
	ID                  int       `json:"id" gorm:"primaryKey"`
	Username            string    `json:"username"`
	Level               int       `json:"level"`
	XP                  int       `json:"xp"`
	PokemonCaught       []Pokemon `json:"pokemon_caught" gorm:"many2many:user_pokemons"`
	ChallengesCompleted int       `json:"challenges_completed"`
	DistanceTraveled    float64   `json:"distance_traveled"`
	LastLocation        Location  `json:"last_location" gorm:"embedded"`
}

type Challenge struct {
	ID          int      `json:"id" gorm:"primaryKey"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	PokemonID   int      `json:"pokemon_id"`
	Pokemon     Pokemon  `json:"pokemon" gorm:"foreignKey:PokemonID"`
	Completed   bool     `json:"completed"`
	Location    Location `json:"location" gorm:"embedded"`
}

func fetchPokemonFromAPI(pokemonID int) (Pokemon, error) {
	url := fmt.Sprintf("%s%d", pokeapiURL, pokemonID)

	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("PokéAPI error: %v", resp.Status)
	}

	var apiData struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Sprites struct {
			FrontDefault string `json:"front_default"`
		} `json:"sprites"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiData); err != nil {
		return Pokemon{}, err
	}

	return Pokemon{
		ID:     apiData.ID,
		Name:   apiData.Name,
		Type:   apiData.Types[0].Type.Name,
		Sprite: apiData.Sprites.FrontDefault,
	}, nil
}

func selectRandomPokemon() (Pokemon, error) {
	randomID := rand.Intn(151) + 1
	return fetchPokemonFromAPI(randomID)
}

func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371
	lat1Rad := lat1 * math.Pi / 180
	lng1Rad := lng1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lng2Rad := lng2 * math.Pi / 180
	deltaLat := lat2Rad - lat1Rad
	deltaLng := lng2Rad - lng1Rad
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func generateRandomLocation(centerLat, centerLng, radiusKm float64) Location {
	angle := rand.Float64() * 2 * math.Pi
	radius := rand.Float64() * radiusKm
	latOffset := radius * math.Cos(angle) / 111.32
	lngOffset := radius * math.Sin(angle) / (111.32 * math.Cos(centerLat*math.Pi/180))
	return Location{Lat: centerLat + latOffset, Lng: centerLng + lngOffset}
}
