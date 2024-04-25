package main

import (
	"log"
	"net"
	"os"

	internal "github.com/HackJack14/SteamSync/Client/Internal"
	window "github.com/HackJack14/SteamSync/Client/Window"
)

func main() {
	clientApp := window.NewClientApp()
	clientApp.Init()
	clientApp.OkClicked = start
	clientApp.CancelClicked = exit
	clientApp.Show()

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

func start(cla *window.ClientApp) {
	ipaddress := cla.GetIP()
	request := cla.GetRequest()
	game := cla.GetGame()
	dir := cla.GetDir()

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

func exit(cla *window.ClientApp) {
	cla.Close()
	os.Exit(0)
}
