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

		request, err := networking.DeserializeString(conn)
		if err != nil {
			fmt.Println(err.Error() + ": connection closed gracefully")
			continue
		}

		switch request.Actstr {
		case "UPLOAD":
			upload(conn)
		case "DOWNLOAD":
			download(conn)
		}
	}
}

func upload(conn net.Conn) {
	defer conn.Close()
	fmt.Println("UPLOAD")

	game, err := networking.DeserializeString(conn)
	if err != nil {
		fmt.Println(err.Error() + ": connection closed gracefully")
		return
	}
	_, err = os.Stat(game.Actstr)
	if os.IsNotExist(err) {
		err = os.Mkdir(game.Actstr, 0755)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	numfiles := networking.DeserializeInt(conn)
	for i := 0; i < int(numfiles); i++ {
		netfile, err := networking.DeserializeFile(conn)
		if err != nil {
			fmt.Println(err.Error() + ": connection closed gracefully")
			return
		}
		file, err := os.Create(filepath.Join(".", game.Actstr, netfile.Name.Actstr))
		if err != nil {
			panic(err)
		}
		io.Copy(file, netfile.Actfile)
		file.Close()
	}
}

func download(conn net.Conn) {
	defer conn.Close()
	fmt.Println("DOWNLOAD")

	game, err := networking.DeserializeString(conn)
	if err != nil {
		fmt.Println(err.Error() + ": connection closed gracefully")
		return
	}
	dir := filepath.Join(".", game.Actstr)
	_, err = os.Stat(dir)
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
