const pokedexGrid = document.getElementById("pokedex-grid");
const caughtCount = document.getElementById("caught-count");
const userId = localStorage.getItem("userId");

const typeColors = {
  fire: "#ff6b6b",
  water: "#1e90ff",
  grass: "#2ed573",
  electric: "#feca57",
  psychic: "#a29bfe",
  ice: "#74b9ff",
  dragon: "#6c5ce7",
  dark: "#636e72",
  fairy: "#fab1a0",
  normal: "#dfe6e9",
  fighting: "#d63031",
  flying: "#81ecec",
  poison: "#a29bfe",
  ground: "#e17055",
  rock: "#b2bec3",
  bug: "#55efc4",
  ghost: "#636e72",
  steel: "#b2bec3"
};

async function loadPokedex() {
  try {
    const res = await fetch(`https://pokeconquest.onrender.com/api/users/${userId}`);
    if (!res.ok) throw new Error("Failed to fetch user data");

    const userData = await res.json();
    const caughtPokemon = userData.pokemon_caught;
    caughtCount.innerText = `You caught ${caughtPokemon.length} Pokémon!`;

 
    for (const pokemon of caughtPokemon) {
      const pokeData = await fetchPokeDetails(pokemon.name.toLowerCase());
      createPokemonCard(pokeData);
    }
  } catch (err) {
    console.error("Error loading Pokédex:", err);
    caughtCount.innerText = "⚠️ Failed to load Pokédex.";
  }
}


async function fetchPokeDetails(name) {
  try {
    const res = await fetch(`https://pokeapi.co/api/v2/pokemon/${name}`);
    if (!res.ok) throw new Error(`Failed to fetch ${name} from PokéAPI`);
    const data = await res.json();


    return {
      id: data.id,
      name: capitalize(data.name),
      sprite: data.sprites.front_default,
      types: data.types.map(t => t.type.name),
      abilities: data.abilities.map(a => a.ability.name),
      stats: {
        hp: data.stats.find(s => s.stat.name === "hp").base_stat,
        attack: data.stats.find(s => s.stat.name === "attack").base_stat,
        defense: data.stats.find(s => s.stat.name === "defense").base_stat
      },
      height: data.height / 10, 
      weight: data.weight / 10  
    };
  } catch (err) {
    console.error(err);
    return {
      id: "???",
      name: capitalize(name),
      sprite: "images/missingno.png", 
      types: ["unknown"],
      abilities: ["N/A"],
      stats: { hp: 0, attack: 0, defense: 0 },
      height: "N/A",
      weight: "N/A"
    };
  }
}

function createPokemonCard(pokemon) {
  const card = document.createElement("div");
  card.className = "pokemon-card";


  const typeColor = typeColors[pokemon.types[0]] || "#b2bec3";

 
  const typeBadges = pokemon.types.map(type => `
    <span class="pokemon-type" style="background-color:${typeColors[type] || '#ccc'}">
      ${capitalize(type)}
    </span>`).join(" ");

  card.innerHTML = `
    <h3>#${pokemon.id} ${pokemon.name}</h3>
    <img src="${pokemon.sprite}" alt="${pokemon.name}" />
    <div>${typeBadges}</div>
    <div class="pokemon-stats">
      <div class="stat">
        <span>HP:</span>
        <div class="stat-bar">
          <div class="stat-fill" style="width:${pokemon.stats.hp}%; background:#ff6b6b;"></div>
        </div>
        <span>${pokemon.stats.hp}</span>
      </div>
      <div class="stat">
        <span>Attack:</span>
        <div class="stat-bar">
          <div class="stat-fill" style="width:${pokemon.stats.attack}%; background:#feca57;"></div>
        </div>
        <span>${pokemon.stats.attack}</span>
      </div>
      <div class="stat">
        <span>Defense:</span>
        <div class="stat-bar">
          <div class="stat-fill" style="width:${pokemon.stats.defense}%; background:#1e90ff;"></div>
        </div>
        <span>${pokemon.stats.defense}</span>
      </div>
      <p><strong>Abilities:</strong> ${pokemon.abilities.join(", ")}</p>
      <p><strong>Height:</strong> ${pokemon.height}m</p>
      <p><strong>Weight:</strong> ${pokemon.weight}kg</p>
    </div>
  `;
  pokedexGrid.appendChild(card);
}

function capitalize(str) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

function filterPokemon() {
  const searchTerm = document.getElementById("search-bar").value.toLowerCase();
  const cards = pokedexGrid.querySelectorAll(".pokemon-card");
  cards.forEach(card => {
    const name = card.querySelector("h3").innerText.toLowerCase();
    card.style.display = name.includes(searchTerm) ? "block" : "none";
  });
}

function goToMap() {
  window.location.href = "map.html";
}


loadPokedex();
