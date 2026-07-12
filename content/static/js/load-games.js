async function getGames() {
  const response = await fetch("/api/getGames");
  var games = await response.json();
  return games;
}

const list = document.getElementById("game-list");
function displayGame(game) {
  var element = document.createElement("li");
  element.textContent = game;
  list.appendChild(element);
}

async function loadGames() {
  var games = await getGames();
  games.forEach(displayGame);
}

loadGames();
