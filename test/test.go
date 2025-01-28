package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":11111")
	if err != nil {
		log.Fatalln(err)
	}

	listener.Accept()
}
