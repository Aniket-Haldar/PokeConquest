# 🌏 Pokémon Adventure - Location-Based Catching Game

An interactive Pokémon adventure game where players explore a real-world map, discover nearby Pokémon, and catch them with exciting animations. Built with **Go (Golang)** backend, **Leaflet.js** for maps, and public APIs for weather and location awareness.  

---

## 🚀 Features
- 🗺 **Map-Based Exploration** – Navigate a real-world map to find Pokémon.  
- ⚡ **Animated Catching Mechanism** – Interactive animations for catching Pokémon.  
- 📖 **Pokédex** – View all Pokémon you have caught with sprites and types.  
- 🌦 **Live Weather Integration** – Weather affects Pokémon spawns (via OpenWeather API).  
- 📍 **Location Awareness** – Shows player’s city and Points of Interest (POIs).  
- 🏆 **Challenges** – Location-based tasks to earn XP and level up.  
- 📊 **Progress Tracking** – XP bar, level, and distance traveled.  

---

## 🌐 Public APIs Used
### 1️⃣ [PokéAPI](https://pokeapi.co/)  
- Fetches Pokémon data (name, type, sprite images).  
- Example: `https://pokeapi.co/api/v2/pokemon/25` returns data for Pikachu.  

### 2️⃣ [OpenWeatherMap API](https://openweathermap.org/api)  
- Retrieves live weather for player’s location.  
- Example: `https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API_KEY}`  

### 3️⃣ [Nominatim API (OpenStreetMap)](https://nominatim.org/release-docs/develop/api/Search/)  
- Converts GPS coordinates to human-readable addresses.  
- Example: Reverse geocoding to get city name and landmarks.  

---

## 📂 Project Structure

├── backend/ # Golang backend
│ ├── main.go # Entry point
│ ├── models.go # Models (User, Pokémon, etc.)
│ ├── controllers.go # API handlers
│ ├── database.go # DB initialization (PostgreSQL)
│ ├── seed.go # Initial data seeding
├── frontend/ # HTML/CSS/JS files
│ ├── index.html # Main map UI
│ ├── css/
│ ├── js/
│ ├── pokedex.html # Pokédex view
├── README.md
├── go.mod / go.sum # Golang dependencies

1. **Clone the repository**
   ```bash
   git clone https://github.com/Aniket-Haldar/pokemon-adventure.git
   cd pokemon-adventure/backend
2. **Configure PostgreSQL**

Create a database called pokequest.

Update database.go with your DB credentials:

dsn := "host=localhost user=postgres password=postgres dbname=pokequest port=5432 sslmode=disable"
3. **Run the backend server**

go mod tidy
go run main.go
✅ Server starts on http://localhost:8080


Endpoint	Method	Description
/api/users	POST	Create a new user
/api/users/{id}	GET	Get user profile
/api/pokemon/nearby?lat&lng	GET	Get nearby Pokémon spawns
/api/catch?user_id&pokemon_id	POST	Catch a Pokémon
/api/weather?lat&lng	GET	Get weather data for current location