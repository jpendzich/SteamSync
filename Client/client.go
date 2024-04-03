package main

import (
	"fmt"
	"net"
	"os"

	networking "github.com/HackJack14/SteamSync/Networking"
)

func main() {
	request := os.Args[1]
	dir := os.Args[2]
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	networking.SerializeString(networking.BuildNetstring(request), conn)

	networking.SerializeInt(uint64(len(names)), conn)
	for _, name := range names {
		fmt.Println(name)
		file, err := os.Open(dir + "/" + name)
		if err != nil {
			panic(err)
		}
		networking.SerializeFile(networking.BuildNetfile(file), conn)
	}
}
