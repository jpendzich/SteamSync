package main

import (
	"fmt"
	"io"
	"net"
	"os"

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

		utils.ReadFile()
	}
}
