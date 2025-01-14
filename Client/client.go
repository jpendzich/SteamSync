package main

import (
	"log"
	"net"
	"os"

	"github.com/schollz/peerdiscovery"
)

var peerAddress string

func main() {
	discoveries, err := peerdiscovery.Discover(peerdiscovery.Settings{Limit: 1})
	if err != nil {
		log.Fatalln(err)
	}
	if len(discoveries) < 1 {
		log.Fatalln("didnt find any connections")
	}
	peerAddress = discoveries[0].Address

	if os.Args[1] == "1" {
		listen()
	} else {
		send()
	}
}

func listen() {
	listener, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(buf))
}

func send() {
	conn, err := net.Dial("tcp", net.JoinHostPort(peerAddress, "10000"))
	if err != nil {
		log.Fatalln(err)
	}
	buf := []byte("das ist ein test string")
	conn.Write(buf)
}
