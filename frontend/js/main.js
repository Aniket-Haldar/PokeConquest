let userId = null;

function startGame() {
  const username = document.getElementById('username').value.trim();
  if (!username) {
    alert("Please enter your trainer name.");
    return;
  }

  fetch("https://pokeconquest.onrender.com//api/users", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username })
  })
    .then(res => res.json())
    .then(data => {
      userId = data.id;
      localStorage.setItem("userId", userId);
      localStorage.setItem("username", data.username);
      alert(`Welcome, ${data.username}!`);
      window.location.href = "map.html";
    })
    .catch(err => {
      console.error(err);
      alert("Failed to create user.");
    });
}
