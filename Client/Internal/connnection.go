package internal

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	networking "github.com/HackJack14/SteamSync/Networking"
)

func Upload(conn net.Conn, game string, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		networking.WriteError(true, conn)
		log.Fatalln(err)
	}

	var saves []string
	for _, entry := range entries {
		if !entry.IsDir() {
			saves = append(saves, entry.Name())
		}
	}

	networking.WriteString(networking.BuildNetstring("UPLOAD"), conn)
	networking.WriteString(networking.BuildNetstring(game), conn)
	networking.WriteInt(uint64(len(saves)), conn)
	for _, save := range saves {
		fmt.Println(save)
		file, err := os.Open(dir + "/" + save)
		if err != nil {
			networking.WriteError(true, conn)
			log.Fatalln(err)
		}
		netfile := networking.BuildNetfile(file)
		networking.WriteFile(netfile, conn)
	}
}

func Download(conn net.Conn, game string, dir string) {
	fmt.Println(game + dir)
	networking.WriteString(networking.BuildNetstring("DOWNLOAD"), conn)
	networking.WriteString(networking.BuildNetstring(game), conn)

	numsaves, err := networking.ReadInt(conn)
	if err != nil {
		log.Fatalf("%s: connection closed gracefully\n", err)
	}
	for i := 0; i < int(numsaves); i++ {
		netfile, err := networking.ReadFile(conn)
		if err != nil {
			log.Fatalf("%s: connection closed gracefully\n", err)
		}
		file, err := os.Create(filepath.Join(dir, netfile.Name.Actstr))
		if err != nil {
			networking.WriteError(true, conn)
			log.Fatalln(err)
		}
		io.Copy(file, netfile.Actfile)
		file.Close()
	}
}
