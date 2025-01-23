package main

import (
	"log"

	"github.com/HackJack14/SteamSync/network"
	"github.com/HackJack14/SteamSync/window"
)

func main() {
	peer, err := network.GetAllPeers()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(peer)
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
