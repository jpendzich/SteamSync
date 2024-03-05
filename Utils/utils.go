package utils

import (
	"fmt"
	"io"
	"net"
)

func ReadString(conn net.Conn) string {
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

func ReadFile(conn net.Conn) {
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
