package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		readFile(conn)
	}
}

func readString(conn net.Conn) string {
	defer conn.Close()
	var buffers []byte

	buffer := make([]byte, 1024)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error:", err)
			break
		}

		buffers = append(buffers, buffer[:length]...)
	}
	return string(buffers)
}

func readFile(conn net.Conn) {
	var fileBytes []byte

	buffer := make([]byte, 1024)
	for {
		len, err := conn.Read(buffer)
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

	file, err := os.Create("README1.md")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	len, err := file.Write(fileBytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len)
}
