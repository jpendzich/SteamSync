package main

import (
	"log"
	"net"
	"os"

	internal "github.com/HackJack14/SteamSync/Client/Internal"
	window "github.com/HackJack14/SteamSync/Client/Window"
)

func main() {
	window.OkClicked = start
	window.CancelClicked = exit

	window.CreateWindow()

	// if len(os.Args) == 1 {
	// 	fmt.Println("Commands:")
	// 	fmt.Println("\t<IPAddress> UPLOAD <Game> <Directory with Savefiles>")
	// 	fmt.Println("\t<IPAddress> DOWNLOAD <Game> <Where to save the Savefiles>")
	// }

	// ipaddress := os.Args[1]
	// request := os.Args[2]
	// game := os.Args[3]
	// dir := os.Args[4]

}

func start() {
	ipaddress := window.EntryIP.Text
	request := window.EntryRequest.Text
	game := window.EntryGame.Text
	dir := window.EntryDir.Text

	conn, err := net.Dial("tcp", ipaddress+":8080")
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

func exit() {
	window.CloseWindow()
	os.Exit(0)
}
