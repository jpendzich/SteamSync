package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/HackJack14/SteamSync/internal"
	"github.com/HackJack14/SteamSync/network"
	"github.com/HackJack14/SteamSync/window"
)

var currentPeer network.Peer

func main() {
	network.RegisterHandler()
	go HandleNewConnections()
	window.OnDownloadGame = DownloadGame
	window.OnUploadGame = UploadGame
	window.OnSelectedPeer = SelectedPeer
	window.OnSelectedGame = SelectedGame
	window.RenderWindow()
}

func UploadGame(game string) error {
	log.Println("uploading game")
	conn, err := net.Dial("tcp", net.JoinHostPort(currentPeer.IpAdress, "9998"))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("dialed")
	connection := network.NewConnection(&conn)
	sender := network.NewRequestSender(connection)
	entries, err := os.ReadDir(game)
	if err != nil {
		log.Fatalln(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			request := network.NewUploadFileRequest()
			request.Game = game
			request.Save = entry.Name()
			response := network.NewUploadFileResponse()
			file, err := os.Create(filepath.Join(game, entry.Name()))
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
			info, err := file.Stat()
			if err != nil {
				log.Fatalln(err)
			}
			sender.SendRequestWriteBinary(1, request, response, file, uint64(info.Size()))
		}
	}
	return nil
}

func DownloadGame(game string) error {
	return nil
}

func SelectedPeer(peer network.Peer) {
	currentPeer = peer
}

func HandleNewConnections() {
	listener, err := net.Listen("tcp", ":9998")
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	connection := network.NewConnection(&conn)

	for {
		packetType, err := connection.ReadPacketType()
		if err != nil {
			log.Fatalln(err)
		}
		err = network.HandlePacket(connection, packetType)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func SelectedGame(game string) error {
	access := internal.NewDbAccess()
	defer access.Close()

	log.Println(game)
	save := access.GetSaveEntryByName(game)
	log.Println(save.Game)
	saveLocations := make(map[string]interface{})
	err := json.Unmarshal([]byte(save.SaveLocation), &saveLocations)
	saveLocation := saveLocations["Windows"].(string)
	fullSaveLocation, err := internal.ConvertToRealPath(save.AppId, saveLocation)
	if err != nil {
		return err
	}
	_, err = os.Stat(fullSaveLocation)
	if err != nil {
		return err
	}
	log.Println(fullSaveLocation)
	return nil
}
