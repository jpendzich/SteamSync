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
		log.Println("awaiting connection")
		conn, err := listener.Accept()
		log.Printf("accepted connection: %s\n", conn.RemoteAddr())
		if err != nil {
			log.Println(err)
			continue
		}

		request, err := networking.ReadString(conn)
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
	log.Println("started receiving upload")

	game, err := networking.ReadString(conn)
	if err != nil {
		log.Printf("%s: connection closed\n", err)
		return
	}
	_, err = os.Stat(game.Actstr)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(game.Actstr, 0755)
		if err != nil {
			log.Println(err)
			return
		}
	} else if err != nil {
		log.Println(err)
		return
	}
	numfiles, err := networking.ReadInt(conn)
	if err != nil {
		log.Fatalf("%s: connection closed\n", err)
	}
	for i := 0; i < int(numfiles); i++ {
		netfile, err := networking.ReadFile(conn)
		if err != nil {
			log.Printf("%s: connection closed gracefully\n", err)
			return
		}
		file, err := os.Create(filepath.Join(".", game.Actstr, netfile.Name.Actstr))
		if err != nil {
			log.Println(err)
			return
		}
		io.Copy(file, netfile.Actfile)
		file.Close()
	}
	log.Println("stopped receiving upload")
}

func download(conn net.Conn) {
	defer conn.Close()
	log.Println("started providing download")

	game, err := networking.ReadString(conn)
	if err != nil {
		log.Println(err)
		return
	}
	dir := filepath.Join(".", game.Actstr)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		log.Println(err)
		return
	}

	names, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}
	networking.WriteInt(uint64(len(names)), conn)
	for _, name := range names {
		if !name.IsDir() {
			file, err := os.Open(filepath.Join(dir, name.Name()))
			if err != nil {
				log.Println(err)
				return
			}
			netfile := networking.BuildNetfile(file)
			networking.WriteFile(netfile, conn)
			file.Close()
		}
	}
	log.Println("stopped providing download")
}
