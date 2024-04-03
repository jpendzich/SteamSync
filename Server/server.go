package main

import (
	"fmt"
	"io"
	"net"
	"os"

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

	numfiles := networking.DeserializeInt(conn)
	for i := 0; i < int(numfiles); i++ {
		netfile := networking.DeserializeFile(conn)
		file, err := os.Create("./" + netfile.Name.Actstr)
		if err != nil {
			panic(err)
		}
		io.Copy(file, netfile.Actfile)
		file.Close()
	}
}

func download(conn net.Conn) {
	fmt.Println("DOWNLOAD")
}
