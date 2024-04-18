package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	networking "github.com/HackJack14/SteamSync/Networking"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Commands:")
		fmt.Println("\t<IP-Address>")
		return
	}
	ipaddress := os.Args[1]
	listener, err := net.Listen("tcp", ipaddress+":8080")
	if err != nil {
		log.Fatalf("%s: connection closed", err)
		return
	}
	defer listener.Close()

	log.Printf("server started listening: %s:8080\n", ipaddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("%s: failed accepting connection", err)
			continue
		}

		request, err := networking.DeserializeString(conn)
		if err != nil {
			log.Printf("%s: connection closed gracefully\n", err.Error())
			continue
		}

		switch request.Actstr {
		case "UPLOAD":
			upload(conn)
		case "DOWNLOAD":
			download(conn)
		}
	}
}

func upload(conn net.Conn) {
	defer conn.Close()
	log.Println("starting upload")

	game, err := networking.DeserializeString(conn)
	if err != nil {
		log.Printf("%s: connection closed\n", err)
		return
	}
	_, err = os.Stat(game.Actstr)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(game.Actstr, 0755)
		if err != nil {
			log.Printf("%s: failed to create direcotry\n", err)
			return
		}
	} else if err != nil {
		log.Printf("%s: failed to retrieve directory info\n", err)
		return
	}
	numfiles := networking.DeserializeInt(conn)
	for i := 0; i < int(numfiles); i++ {
		netfile, err := networking.DeserializeFile(conn)
		if err != nil {
			log.Printf("%s: connection closed gracefully\n", err)
			return
		}
		file, err := os.Create(filepath.Join(".", game.Actstr, netfile.Name.Actstr))
		if err != nil {
			log.Printf("%s: failed to create file\n", err)
			return
		}
		io.Copy(file, netfile.Actfile)
		file.Close()
	}
}

func download(conn net.Conn) {
	defer conn.Close()
	log.Println("DOWNLOAD")

	game, err := networking.DeserializeString(conn)
	if err != nil {
		log.Printf("%s: connection closed\n", err)
		return
	}
	dir := filepath.Join(".", game.Actstr)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		log.Printf("%s: directory does not exist\n", err)
		return
	}

	names, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("%s: failed to read directory\n", err)
		return
	}
	networking.SerializeInt(uint64(len(names)), conn)
	for _, name := range names {
		if !name.IsDir() {
			file, err := os.Open(filepath.Join(dir, name.Name()))
			if err != nil {
				log.Printf("%s: failed to open file\n", err)
				return
			}
			netfile := networking.BuildNetfile(file)
			networking.SerializeFile(netfile, conn)
			file.Close()
		}
	}
}
