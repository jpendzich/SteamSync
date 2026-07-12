var searchInput = document.getElementById("search-input");
searchInput.addEventListener("input", searchInputChanged);

var timeoutId = 0;
function searchInputChanged(event) {
  clearTimeout(timeoutId);
  timeoutId = setTimeout(() => {
    loadGames(searchInput.value);
  }, 300);
}
