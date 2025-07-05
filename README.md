# PokéConquest 🌍🎮

A location-based scavenger hunt game where players explore the real world to catch Pokémon, complete challenges, and interact with nearby trainers — inspired by Pokémon Go!

Built with **Go (Golang)** for the backend, **Vanilla HTML/CSS/JS** for the frontend, and **PostgreSQL** as the database.

---

## 🚀 Features
- 📍 **Location-aware gameplay** using browser geolocation
- 🎯 Catch Pokémon around you
- 🧑‍🤝‍🧑 View and interact with nearby trainers
- 🌱 Level up & track XP
- 🧠 Professor Oak’s AI tips (powered by OpenAI API)
- 🌤 Dynamic weather system (OpenWeather API)

---

## 🛠 Tech Stack

| Component          | Technology        |
|---------------------|--------------------|
| Backend API         | Go (Golang)       |
| Database            | PostgreSQL        |
| Frontend            | HTML/CSS/JS       |
| Mapping             | Leaflet.js + OpenStreetMap |
| Weather API         | OpenWeather API   |
| AI Tips             | OpenAI API        |
| Hosting             | Render (Dockerized) |

---

Visit: https://pokeconquest.onrender.com


## Appendix


│
├── frontend/ # Frontend static files
│ ├── index.html # Landing page
│ ├── map.html # Map interface
│ ├── js/ # JavaScript files
│ │ ├── map.js # Map & gameplay logic
│ │ └── ...
│ ├── css/ # Stylesheets
│ └── images/ # Game assets
│
├── go.mod # Go module definition
├── main.go # Backend entrypoint
├── database.go # Database connection & models
├── controllers.go # API route handlers
├── Dockerfile # Docker build config
├── README.md # You are here 😉
└── .env.example # Sample environment variables
---

## ⚙️ API Endpoints

### 👤 User
| Method | Endpoint                          | Description                       |
|--------|-------------------------------------|-------------------------------------|
| POST   | `/api/users`                      | Create a new user                  |
| GET    | `/api/users/{id}`                 | Get user profile                   |
| PUT    | `/api/users/{id}/location`        | Update user’s location             |
| GET    | `/api/users/{id}/gamestate`       | Get user’s game state              |

### 🗺 Pokémon & Challenges
| Method | Endpoint                          | Description                       |
|--------|-------------------------------------|-------------------------------------|
| GET    | `/api/pokemon/nearby?lat=&lng=`   | Get Pokémon near user’s location   |
| POST   | `/api/catch?user_id=&pokemon_id=` | Attempt to catch a Pokémon         |
| GET    | `/api/challenges`                 | Get available challenges           |
| POST   | `/api/challenges/{id}/complete`   | Mark challenge as complete         |

### 🌤 Environment
| Method | Endpoint          | Description                 |
|--------|---------------------|-----------------------------|
| GET    | `/api/health`     | Health check (for Render)   |
| GET    | `/api/locations/pois` | Get nearby points of interest |

### 🧠 AI Tips
| Method | Endpoint          | Description                 |
|--------|---------------------|-----------------------------|
| POST   | `/api/ai/tip`     | Get a gameplay tip           |
| POST   | `/api/ai/strategy`| Get a strategy suggestion    |

---


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

DATABASE_URL=postgres://pokeuser:pokepass@localhost:5432/pokequest
OPENWEATHER_API_KEY=your_openweather_api_key
OPENAI_API_KEY=your_openai_api_key

## 🖥 Local Development

### Prerequisites
- [Go](https://golang.org/dl/) >= 1.19
- [PostgreSQL](https://www.postgresql.org/) >= 13
- [Node.js](https://nodejs.org/) (for frontend tweaks, optional)

---

### ⚡ Install & Run

1️⃣ Clone the repo:
```bash
git clone https://github.com/your-username/pokeconquest.git
cd pokeconquest


cp .env.example .env


DATABASE_URL=postgres://pokeuser:pokepass@localhost:5432/pokequest
OPENWEATHER_API_KEY=your_openweather_api_key
OPENAI_API_KEY=your_openai_api_key


createdb pokequest

go run main.go

 Open http://localhost:8080 in your browser.


 docker build -t pokeconquest .
docker run -p 8080:8080 pokeconquest

```