package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ApiGamesHandler struct {
	games []string
}

func NewApiGamesHandler(games []string) *ApiGamesHandler {
	return &ApiGamesHandler{
		games: games,
	}
}

func (h *ApiGamesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")
	results := make([]string, 0, 100)
	var manifestJson []byte
	if search != "" {
		gamesFound := 0
		for _, game := range h.games {
			if gamesFound > 100 {
				break
			}

			if strings.Contains(game, search) {
				gamesFound++
				results = append(results, game)
			}
		}
	} else {
		results = append(results, h.games[:100]...)
	}

	manifestJson, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	w.Write(manifestJson)
}
