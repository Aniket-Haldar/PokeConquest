const userId = localStorage.getItem("userId");

function loadPokedex() {
  fetch(`http://localhost:8080/api/users/${userId}`)
    .then(res => res.json())
    .then(data => {
      const pokedex = document.getElementById("pokedex");
      if (data.pokemon_caught.length === 0) {
        pokedex.innerHTML = "<p>No Pokémon caught yet!</p>";
        return;
      }
      data.pokemon_caught.forEach(p => {
        const card = document.createElement("div");
        card.innerHTML = `
          <h3>${p.name}</h3>
          <img src="${p.sprite}" width="80" />
          <p>Type: ${p.type}</p>
        `;
        pokedex.appendChild(card);
      });
    })
    .catch(err => console.error("Failed to load Pokédex:", err));
}

function goBack() {
  window.location.href = "map.html";
}

window.onload = loadPokedex;
