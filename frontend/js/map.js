const userId = localStorage.getItem("userId");
if (!userId) {
  alert("⚠️ No user found. Please start from the landing page.");
  window.location.href = "index.html";
}

let map;
let factInterval;

const pokemonFacts = [
  "🐱 Meowth is the only Pokémon that can speak human language in the anime.",
  "🔥 Charmander’s flame goes out if it dies.",
  "⚡ Pikachu stores electricity in its cheeks.",
  "🌿 Bulbasaur was the first Pokémon ever designed.",
  "🦇 Zubat has no eyes but uses echolocation.",
  "🥚 Togepi’s shell acts as an energy absorber.",
  "🥊 Hitmonlee and Hitmonchan are named after Bruce Lee & Jackie Chan.",
  "🌊 Magikarp was inspired by a Chinese legend of a carp turning into a dragon.",
  "🧠 Alakazam has an IQ of over 5000.",
  "🌙 Clefairy was almost the series mascot instead of Pikachu!"
];

function initMap(lat, lng) {
  const darkMode = document.body.classList.contains('dark-mode');
  const tileUrl = darkMode
    ? 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png'
    : 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';

  map = L.map('map').setView([lat, lng], 15);

  L.tileLayer(tileUrl, {
    attribution: '&copy; OpenStreetMap contributors'
  }).addTo(map);

  L.marker([lat, lng]).addTo(map)
    .bindPopup("📍 You are here").openPopup();

  fetchNearbyPokemon(lat, lng);
  fetchNearbyTrainers(lat, lng);
  fetchWeather(lat, lng);
  loadTrainerProfile();

  hideLoader();
}

function fetchNearbyPokemon(lat, lng) {
  showLoader();
  fetch(`https://pokeconquest.onrender.com/api/pokemon/nearby?lat=${lat}&lng=${lng}`)
    .then(res => res.json())
    .then(data => {
      if (Array.isArray(data)) {
        data.forEach(spawn => addPokemonMarker(spawn.pokemon, spawn.location));
      } else {
        console.error("Unexpected Pokémon data format:", data);
        showError("⚠️ MissingNo appeared! Could not load Pokémon.");
      }
    })
    .catch(err => {
      console.error("Failed to load Pokémon:", err);
      showError("⚠️ MissingNo appeared! Could not fetch Pokémon.");
    })
    .finally(hideLoader);
}

function fetchNearbyTrainers(lat, lng) {
  fetch(`https://pokeconquest.onrender.com/api/trainers/nearby?lat=${lat}&lng=${lng}`)
    .then(res => res.json())
    .then(trainers => {
      trainers.forEach(trainer => {
        const trainerIcon = L.icon({
          iconUrl: "images/trainer.png",
          iconSize: [40, 40],
          iconAnchor: [20, 40],
          popupAnchor: [0, -35]
        });
        const marker = L.marker(
          [trainer.last_location.lat, trainer.last_location.lng],
          { icon: trainerIcon }
        ).addTo(map);
        marker.bindPopup(`
          <div style="text-align:center;">
            <h4>${trainer.username}</h4>
            <p>Level ${trainer.level}</p>
          </div>
        `);
      });
    })
    .catch(err => console.error("Failed to load trainers:", err));
}

function fetchWeather(lat, lng) {
  const apiKey = "04fb8c6c637808d39017edf779f7bafd";
  fetch(`https://api.openweathermap.org/data/2.5/weather?lat=${lat}&lon=${lng}&units=metric&appid=${apiKey}`)
    .then(res => res.json())
    .then(data => {
      const weatherMain = data.weather[0].main;
      const description = data.weather[0].description;
      const temp = data.main.temp;
      const isNight = (data.dt < data.sys.sunrise || data.dt > data.sys.sunset);

      document.getElementById("weather").innerText =
        `${isNight ? "🌙 Night" : "🌤 Day"} • ${weatherMain}: ${description}, ${temp}°C`;

      applyWeatherTheme(weatherMain, isNight);
    })
    .catch(err => {
      console.error("Failed to load weather:", err);
      document.getElementById("weather").innerText = "⚠️ Failed to load weather!";
    });
}

function applyWeatherTheme(weather, isNight) {
  document.body.classList.remove("sunny-theme", "rainy-theme", "cloudy-theme", "night-theme");
  if (isNight) {
    document.body.classList.add("night-theme");
  } else {
    switch (weather) {
      case "Clear":
        document.body.classList.add("sunny-theme");
        break;
      case "Rain":
        document.body.classList.add("rainy-theme");
        break;
      case "Clouds":
        document.body.classList.add("cloudy-theme");
        break;
      default:
        document.body.classList.add("sunny-theme"); // fallback
    }
  }
}

function loadTrainerProfile() {
  fetch(`https://pokeconquest.onrender.com/api/users/${userId}`)
    .then(res => res.json())
    .then(user => {
      document.getElementById("trainer-name").innerText = `Trainer: ${user.username}`;
      document.getElementById("trainer-level").innerText = `Level: ${user.level}`;
      const xpPercent = Math.min((user.xp % 100) / 100 * 100, 100);
      document.getElementById("xp-fill").style.width = `${xpPercent}%`;
    })
    .catch(err => {
      console.error("Failed to load trainer profile:", err);
      showError("⚠️ Could not load trainer data.");
    });
}

navigator.geolocation.getCurrentPosition(
  pos => {
    const { latitude, longitude } = pos.coords;
    initMap(latitude, longitude);
    updateUserLocation(latitude, longitude);
  },
  err => {
    alert("⚠️ Please enable location to play PokéQuest.");
    console.error("Geolocation error:", err);
    showError("⚠️ MissingNo appeared! Could not get location.");
    hideLoader();
  }
);

function updateUserLocation(lat, lng) {
  fetch(`https://pokeconquest.onrender.com/api/users/${userId}/location`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ lat, lng })
  }).catch(err => console.error("Failed to update location:", err));
}
function addPokemonMarker(pokemon, location) {
  const marker = L.marker([location.lat, location.lng]).addTo(map);
  const popupContent = `
    <div style="text-align:center;">
      <h4 style="margin:0;color:#2f3542;">${pokemon.name}</h4>
      <img src="${pokemon.sprite}" style="width:60px;height:60px;margin:5px;">
      <br/>
      <button onclick="catchPokemon(${pokemon.id})" style="
        margin-top:5px;
        background:#ff4757;
        color:white;
        border:none;
        padding:6px 12px;
        border-radius:8px;
        cursor:pointer;
        font-weight:bold;
        box-shadow:0 4px 8px rgba(0,0,0,0.3);
      ">🎯 Catch</button>
    </div>
  `;
  marker.bindPopup(popupContent);
}
function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function showLoader() {
  const loader = document.getElementById('loading-overlay');
  const loadingText = document.getElementById('loading-text');
  if (loader && loadingText) {
    loader.classList.remove('hidden');
    const randomFact = pokemonFacts[Math.floor(Math.random() * pokemonFacts.length)];
    loadingText.innerText = randomFact;
    clearInterval(factInterval);
    factInterval = setInterval(() => {
      const newFact = pokemonFacts[Math.floor(Math.random() * pokemonFacts.length)];
      loadingText.innerText = newFact;
    }, 2500);
  }
}
async function catchPokemon(pokemonId) {
  const overlay = document.getElementById('catch-animation');
  const message = document.getElementById('catch-message');

  try {
    overlay.classList.remove('hidden');
    message.textContent = "🎯 Throwing Pokéball...";
    await delay(2000); // Simulate throw

    const response = await fetch(`https://pokeconquest.onrender.com/api/catch?user_id=${userId}&pokemon_id=${pokemonId}`, {
      method: "POST",
    });
    const data = await response.json();

    if (response.ok) {
      const caughtPokemon = data.user.pokemon_caught[data.user.pokemon_caught.length - 1];
      message.textContent = `🎉 Gotcha! You caught ${caughtPokemon.name}!`;
      loadTrainerProfile(); // update XP and level
    } else {
      message.textContent = "❌ Oh no! It escaped!";
    }
    await delay(2000);
  } catch (err) {
    console.error("Failed to catch Pokémon:", err);
    alert("⚠️ Failed to catch Pokémon!");
  } finally {
    overlay.classList.add('hidden');
  }
}

function hideLoader() {
  const loader = document.getElementById('loading-overlay');
  if (loader && !loader.classList.contains('hidden')) {
    loader.classList.add('hidden');
    clearInterval(factInterval);
    factInterval = null;
  }
}

// 🧠 AI TIP System (Professor Oak)
document.getElementById('next-tip').addEventListener('click', fetchAITip);
const aiContainer = document.getElementById('ai-container');
const openAIButton = document.getElementById('open-ai');
const closeAIButton = document.getElementById('close-ai');
closeAIButton.addEventListener('click', () => {
  aiContainer.classList.add('hidden');
  openAIButton.classList.remove('hidden');
});

openAIButton.addEventListener('click', () => {
  aiContainer.classList.remove('hidden');
  openAIButton.classList.add('hidden');
});

const toggleVoiceButton = document.getElementById('toggle-voice');
let oakVoiceEnabled = true; // voice on by default

// Fix voices loading
let voices = [];
window.speechSynthesis.onvoiceschanged = () => {
  voices = speechSynthesis.getVoices();
};

toggleVoiceButton.addEventListener('click', () => {
  oakVoiceEnabled = !oakVoiceEnabled;
  toggleVoiceButton.textContent = oakVoiceEnabled ? "🔊 Mute Oak" : "🔇 Unmute Oak";

  // Stop current speech if muting
  if (!oakVoiceEnabled && speechSynthesis.speaking) {
    speechSynthesis.cancel();
  }
});



function fetchAITip() {
  const aiText = document.getElementById('typed-text');
  aiText.textContent = "Professor Oak is thinking… 🤔";

  fetch('https://pokeconquest.onrender.com/api/ai/tip', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ context: "Give me a helpful Pokémon tip." })
  })
    .then(response => {
      if (!response.ok) throw new Error("Server returned an error");
      return response.json();
    })
    .then(data => typeWriter(data.tip))
    .catch(err => {
      console.error(err);
      aiText.textContent = "Hmm… Something went wrong!";
    });
}

function typeWriter(text) {
  const aiText = document.getElementById('typed-text');
  aiText.textContent = '';
  let i = 0;

  function typing() {
    if (i < text.length) {
      aiText.textContent += text.charAt(i);
      i++;
      setTimeout(typing, 40); // typing speed
    } else {
      speak(text); // 🗣 Speak after typing
    }
  }
  typing();
}

function speak(text) {
  if (!oakVoiceEnabled) return; // respect mute
  if ('speechSynthesis' in window) {
    const utterance = new SpeechSynthesisUtterance(text);
    // Use first English voice or fallback
    const voice = voices.find(v => v.lang.startsWith('en')) || null;
    if (voice) utterance.voice = voice;
    utterance.rate = 0.9;
    utterance.pitch = 1.0;
    utterance.volume = 1;
    speechSynthesis.speak(utterance);
  }
}

function showError(msg) {
  alert(msg);
}

function toggleDarkMode() {
  document.body.classList.toggle('dark-mode');
  const newTheme = document.body.classList.contains('dark-mode') ? 'dark' : 'light';
  localStorage.setItem('theme', newTheme);

  if (map) {
    map.eachLayer(layer => map.removeLayer(layer));
    const tileUrl = newTheme === 'dark'
      ? 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png'
      : 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
    L.tileLayer(tileUrl, {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);
    fetchNearbyPokemon(map.getCenter().lat, map.getCenter().lng);
    fetchNearbyTrainers(map.getCenter().lat, map.getCenter().lng);
  }
}

window.addEventListener('load', () => {
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme === 'dark') {
    document.body.classList.add('dark-mode');
  }

  aiContainer.classList.remove('hidden');
  openAIButton.classList.add('hidden');

  showLoader();
  fetchAITip(); // Ask Oak for advice on load
});
