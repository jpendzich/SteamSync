package main

import (
	"github.com/HackJack14/SteamSync/window"
)

func main() {
	window.OnDownloadGame = DownloadGame
	window.OnUploadGame = UploadGame
	window.RenderWindow()
}

func UploadGame(game string) error {
	return nil
}

func DownloadGame(game string) error {
	return nil
}
