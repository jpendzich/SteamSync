package main

import (
	"fmt"
	"net"
	"strconv"

	utils "github.com/HackJack14/SteamSync/Utils"
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

		request := utils.ReadString(conn)

		switch request {
		case "UPLOAD":
			upload(conn)
		case "DOWNLOAD":
			download(conn)
		}
	}
}

func upload(conn net.Conn) {
	defer conn.Close()
	amountfiles, err := strconv.Atoi(utils.ReadString(conn))
	if err != nil {
		fmt.Println(err)
	}
	game := utils.ReadString(conn)

	for i := 0; i < amountfiles; i++ {
		utils.ReadFile(game, conn)
	}
	fmt.Println("test")
}

func download(conn net.Conn) {
	fmt.Println("DOWNLOAD")
}
