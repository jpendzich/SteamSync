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
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		request := networking.DeserializeString(conn)

		switch request.Actstr {
		case "UPLOAD":
			upload(conn)
		case "DOWNLOAD":
			download(conn)
		}
	}
}

func upload(conn net.Conn) {
	fmt.Println("UPLOAD")

	game := networking.DeserializeString(conn)
	err := os.Mkdir(game.Actstr, 0755)
	if err != nil {
		panic(err)
	}
	numfiles := networking.DeserializeInt(conn)
	for i := 0; i < int(numfiles); i++ {
		netfile := networking.DeserializeFile(conn)
		file, err := os.Create(filepath.Join(".", game.Actstr, netfile.Name.Actstr))
		if err != nil {
			panic(err)
		}
		io.Copy(file, netfile.Actfile)
		file.Close()
	}
}

func download(conn net.Conn) {
	fmt.Println("DOWNLOAD")

	game := networking.DeserializeString(conn)
	dir := filepath.Join(".", game.Actstr)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return
	}

	names, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	networking.SerializeInt(uint64(len(names)), conn)
	for _, name := range names {
		if !name.IsDir() {
			hasError := false
			file, err := os.Open(filepath.Join(dir, name.Name()))
			if err != nil {
				hasError = true
			}
			netfile := networking.BuildNetfile(file)
			netfile.Error = hasError
			networking.SerializeFile(netfile, conn)
			file.Close()
		}
	}
}
