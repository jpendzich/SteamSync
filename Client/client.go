package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	networking "github.com/HackJack14/SteamSync/Networking"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Commands:")
		fmt.Println("\tUPLOAD <Game> <Directory with Savefiles>")
		fmt.Println("\tDOWNLOAD <Game> <Where to save the Savefiles>")
		return
	}

	request := os.Args[1]
	game := os.Args[2]
	dir := os.Args[3]

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	switch request {
	case "UPLOAD":
		upload(conn, game, dir)
	case "DOWNLOAD":
		download(conn, game, dir)
	}

}

func upload(conn net.Conn, game string, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var saves []string
	for _, entry := range entries {
		if !entry.IsDir() {
			saves = append(saves, entry.Name())
		}
	}

	networking.SerializeString(networking.BuildNetstring("UPLOAD"), conn)
	networking.SerializeString(networking.BuildNetstring(game), conn)
	networking.SerializeInt(uint64(len(saves)), conn)
	for _, save := range saves {
		fmt.Println(save)
		hasError := false
		file, err := os.Open(dir + "/" + save)
		if err != nil {
			hasError = true
		}
		netfile := networking.BuildNetfile(file)
		netfile.Error = hasError
		networking.SerializeFile(netfile, conn)
	}
}

func download(conn net.Conn, game string, dir string) {
	fmt.Println(game + dir)
	networking.SerializeString(networking.BuildNetstring("DOWNLOAD"), conn)
	networking.SerializeString(networking.BuildNetstring(game), conn)

	numsaves := networking.DeserializeInt(conn)
	for i := 0; i < int(numsaves); i++ {
		netfile := networking.DeserializeFile(conn)
		file, err := os.Create(filepath.Join(dir, netfile.Name.Actstr))
		if err != nil {
			panic(err)
		}
		io.Copy(file, netfile.Actfile)
		file.Close()
	}
}
