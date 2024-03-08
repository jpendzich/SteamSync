package main

import (
	"fmt"
	"net"
	"os"

	utils "github.com/HackJack14/SteamSync/Utils"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	utils.SendFile(os.Args[1], conn)
}
