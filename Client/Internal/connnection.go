package internal

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	networking "github.com/HackJack14/SteamSync/Networking"
)

func UploadGameSaves(conn net.Conn, game string, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		networking.WriteError(true, conn)
		log.Fatalln(err)
	}

	var saves []string
	for _, entry := range entries {
		if !entry.IsDir() {
			saves = append(saves, entry.Name())
		}
	}

	networking.WriteString(networking.BuildNetstring("UPLOAD"), conn)
	networking.WriteString(networking.BuildNetstring(game), conn)
	networking.WriteInt(uint64(len(saves)), conn)
	for _, save := range saves {
		fmt.Println(save)
		file, err := os.Open(dir + "/" + save)
		if err != nil {
			networking.WriteError(true, conn)
			log.Fatalln(err)
		}
		netfile := networking.BuildNetfile(file)
		networking.WriteFile(netfile, conn)
	}
}

func DownloadGameSaves(conn net.Conn, game string, dir string) {
	fmt.Println(game + dir)
	networking.WriteString(networking.BuildNetstring("DOWNLOAD"), conn)
	networking.WriteString(networking.BuildNetstring(game), conn)

	numsaves, err := networking.ReadInt(conn)
	if err != nil {
		log.Fatalf("%s: connection closed gracefully\n", err)
	}
	for i := 0; i < int(numsaves); i++ {
		netfile, err := networking.ReadFile(conn)
		if err != nil {
			log.Fatalf("%s: connection closed gracefully\n", err)
		}
		file, err := os.Create(filepath.Join(dir, netfile.Name.Actstr))
		if err != nil {
			networking.WriteError(true, conn)
			log.Fatalln(err)
		}
		io.Copy(file, netfile.Actfile)
		file.Close()
	}
}

func GetGames(conn net.Conn) []string {
	networking.WriteString(networking.BuildNetstring("GAMES"), conn)

	numgames, err := networking.ReadInt(conn)
	if err != nil {
		log.Fatalf("%s: connection closed gracfully\n", err)
	}
	games := make([]string, 0)
	for i := 0; i < int(numgames); i++ {
		game, err := networking.ReadString(conn)
		if err != nil {
			log.Fatalf("%s: connection closed gracfully\n", err)
		}
		games = append(games, game.Actstr)
	}
	return games
}
