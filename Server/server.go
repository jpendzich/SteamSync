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
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
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
			upload(conn, exeDir)
		case "DOWNLOAD":
			download(conn, exeDir)
		}
	}
}

func upload(conn net.Conn, exeDir string) {
	defer conn.Close()
	log.Println("started receiving upload")

	game, err := networking.ReadString(conn)
	if err != nil {
		log.Printf("%s: connection closed\n", err)
		return
	}
	dir := filepath.Join(exeDir, game.Actstr)
	_, err = os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			networking.WriteError(true, conn)
			log.Println(err)
			return
		}
	} else if err != nil {
		networking.WriteError(true, conn)
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
			log.Printf("%s: connection closed\n", err)
			return
		}
		file, err := os.Create(filepath.Join(dir, netfile.Name.Actstr))
		if err != nil {
			networking.WriteError(true, conn)
			log.Println(err)
			return
		}
		_, err = io.Copy(file, netfile.Actfile)
		if err != nil {
			networking.WriteError(true, conn)
			log.Println(err)
			return
		}
		file.Close()
	}
	log.Println("stopped receiving upload")
}

func download(conn net.Conn, exeDir string) {
	defer conn.Close()
	log.Println("started providing download")

	game, err := networking.ReadString(conn)
	if err != nil {
		log.Printf("%s: connection closed\n", err)
		return
	}
	dir := filepath.Join(exeDir, game.Actstr)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		networking.WriteError(true, conn)
		log.Println(err)
		return
	}

	names, err := os.ReadDir(dir)
	if err != nil {
		networking.WriteError(true, conn)
		log.Println(err)
		return
	}
	networking.WriteInt(uint64(len(names)), conn)
	for _, name := range names {
		if !name.IsDir() {
			file, err := os.Open(filepath.Join(dir, name.Name()))
			if err != nil {
				networking.WriteError(true, conn)
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
