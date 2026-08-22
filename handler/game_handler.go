package handler

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type GameHandler struct {
	template *template.Template
}

type gameTemplateData struct {
	Title string
}

func NewGameHandler() *GameHandler {
	t, err := template.ParseFiles(filepath.Join(GetExecutableDir(), "dynamic/game/game.html"))
	if err != nil {
		panic(err)
	}
	return &GameHandler{
		template: t,
	}
}

func (h *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := gameTemplateData{
		Title: r.PathValue("name"),
	}
	h.template.Execute(w, data)
}
