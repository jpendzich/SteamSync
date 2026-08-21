package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/goccy/go-yaml"
	"github.com/jpendzich/SteamSync/handler"
)

func main() {
	http.Handle("/", handler.NewStaticHandler())

	manifestFile, err := os.ReadFile(filepath.Join(handler.GetExecutableDir(), "dynamic/manifest/manifest.yaml"))
	if err != nil {
		panic(err)
	}

	manifest := make(map[string]any)
	games := make([]string, 0)

	err = yaml.Unmarshal(manifestFile, &manifest)
	if err != nil {
		panic(err)
	}
	log.Println("Read manifest")

	for k := range manifest {
		games = append(games, k)
	}
	slices.Sort(games)
	log.Println("Built game list")

	http.Handle("/api/games", handler.NewApiGamesHandler(games))

	http.Handle("/api/game", handler.NewApiGameHandler(manifest))

	log.Fatalln(http.ListenAndServe(":8070", nil))
}
