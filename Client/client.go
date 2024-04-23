package main

import (
	"fmt"
	"log"
	"net"
	"os"

	internal "github.com/HackJack14/SteamSync/Client/Internal"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Commands:")
		fmt.Println("\t<IPAddress> UPLOAD <Game> <Directory with Savefiles>")
		fmt.Println("\t<IPAddress> DOWNLOAD <Game> <Where to save the Savefiles>")
	}

	ipadress := os.Args[1]
	request := os.Args[2]
	game := os.Args[3]
	dir := os.Args[4]

	conn, err := net.Dial("tcp", ipadress+":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	switch request {
	case "UPLOAD":
		internal.Upload(conn, game, dir)
	case "DOWNLOAD":
		internal.Download(conn, game, dir)
	}
}
