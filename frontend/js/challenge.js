const userId = localStorage.getItem("userId");

function loadChallenges() {
  fetch(`http://localhost:8080/api/users/${userId}/gamestate`)
    .then(res => res.json())
    .then(data => {
      const challengesDiv = document.getElementById("challenges");
      data.active_challenges.forEach(challenge => {
        const card = document.createElement("div");
        card.innerHTML = `
          <h3>${challenge.title}</h3>
          <p>${challenge.description}</p>
          <button onclick="completeChallenge(${challenge.id})">Complete</button>
        `;
        challengesDiv.appendChild(card);
      });
    })
    .catch(err => console.error("Failed to load challenges:", err));
}

function completeChallenge(challengeId) {
  navigator.geolocation.getCurrentPosition(pos => {
    const { latitude, longitude } = pos.coords;

    fetch(`http://localhost:8080/api/challenges/${challengeId}/complete`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        user_id: parseInt(userId),
        location: { lat: latitude, lng: longitude }
      })
    })
      .then(res => res.json())
      .then(data => {
        alert(`Challenge Complete! You caught ${data.pokemon.name} and gained ${data.xp_gained} XP.`);
        window.location.reload();
      })
      .catch(err => console.error("Failed to complete challenge:", err));
  });
}

function goBack() {
  window.location.href = "map.html";
}

window.onload = loadChallenges;
