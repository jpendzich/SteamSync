package main

import (
	"fmt"
	"net"
	"os"

	utils "github.com/HackJack14/SteamSync/Utils"
)

func main() {
	/*
		dir := os.Args[2]
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
		}

		var names []string
		for _, entry := range entries {
			if !entry.IsDir() {
				names = append(names, entry.Name())
			}
		}
	*/

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	utils.WriteString(os.Args[1], conn)

	/*
		utils.WriteString(strconv.Itoa(len(names)), conn)

		utils.WriteString("Fallout3", conn)

		for _, name := range names {
			fmt.Println(name)
			utils.SendFile(filepath.Join(dir, name), conn)
		}
	*/
}
