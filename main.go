package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

func main() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exePath := filepath.Dir(exe)

	http.Handle("/", http.FileServer(http.Dir(filepath.Join(exePath, "static"))))

	manifestFile, err := os.ReadFile(filepath.Join(exePath, "static/manifest/manifest.yaml"))
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

	for k, _ := range manifest {
		games = append(games, k)
	}
	slices.Sort(games)
	log.Println("Built game list")

	http.HandleFunc("/api/getGames", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		search := query.Get("search")
		results := make([]string, 0, 100)
		var manifestJson []byte
		if search != "" {
			gamesFound := 0
			for _, game := range games {
				if gamesFound > 100 {
					break
				}

				if strings.Contains(game, search) {
					gamesFound++
					results = append(results, game)
				}
			}
		} else {
			results = append(results, games[:100]...)
		}

		manifestJson, err = json.Marshal(results)
		if err != nil {
			panic(err)
		}
		w.Write(manifestJson)
	})

	log.Fatalln(http.ListenAndServe(":8070", nil))
}
