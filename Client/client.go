package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	sendFile(conn)
}

func sendFile(conn net.Conn) {
	file, err := os.Open("README.md")
	if err != nil {
		fmt.Println(err)
	}

	var fileBytes []byte

	buffer := make([]byte, 1024)
	for {
		len, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
			}
		}
		buffer = buffer[:len]
		fileBytes = append(fileBytes, buffer...)
	}

	conn.Write(fileBytes)
}
