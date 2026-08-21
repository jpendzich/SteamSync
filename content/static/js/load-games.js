const list = document.getElementById("game-list");

async function getGames(search) {
    let url = new URL("/api/games", window.location.origin);
    url.searchParams.append("search", search);
    const response = await fetch(url.toString());
    var games = await response.json();
    return games;
}

function displayGame(game) {
    var listElement = document.createElement("li");
    listElement.textContent = game;
    listElement.addEventListener("click", () => {
        let url = new URL("/api/game", window.location.origin)
        url.searchParams.append("name", listElement.innerText)
        window.location.href = url.toString()
    })
    list.appendChild(listElement);
}

async function loadGames(search) {
    while (list.firstChild) {
        list.removeChild(list.firstChild);
    }
    var games = await getGames(search);
    games.forEach(displayGame);
}

loadGames("");
