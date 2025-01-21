package main

import (
	"log"
	"net"
	"os"
	"path"

	"github.com/HackJack14/SteamSync/network"
)

func main() {
	if os.Args[1] == "1" {
		listen()
	} else {
		send()
	}
}

func listen() {
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	connection := network.NewConnection(&conn)
	network.RegisterHandler()
	packetType, err := connection.ReadPacketType()
	if err != nil {
		log.Fatalln(err)
	}
	err = network.HandlePacket(connection, packetType)
	if err != nil {
		log.Fatalln(err)
	}
}

func send() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	connection := network.NewConnection(&conn)
	sender := network.NewRequestSender(connection)

	// file transfer worked successfully
	test := network.NewDownloadFileRequest()
	test.Game = "Fallout 3"
	test.Save = "1-Outside"

	testResponse := network.NewDownloadFileResponse()

	file, err := os.Create(path.Join(test.Game, test.Save))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	log.Println(testResponse)
	err = sender.SendRequestReadBinary(test.PacketType, test, testResponse, file)
	if err != nil {
		log.Fatalln(err)
	}

	if testResponse.ErrorCode == 0 {
		log.Println("successfully downloaded file")
	}
}
