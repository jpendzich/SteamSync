package handler

import (
	"encoding/json"
	"net/http"
)

type ApiGameHandler struct {
	manifest map[string]any
}

func NewApiGameHandler(manifest map[string]any) *ApiGameHandler {
	return &ApiGameHandler{
		manifest: manifest,
	}
}

func (h *ApiGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if game, ok := h.manifest[name]; ok {
		gameJson, err := json.Marshal(game)
		if err != nil {
			panic(err)
		}
		w.Write(gameJson)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
