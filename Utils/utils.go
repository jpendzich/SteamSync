package utils

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
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
	var filebytes []byte

	buffer := make([]byte, 1024)
	len, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	filename := string(buffer[:len])
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
		filebytes = append(filebytes, buffer...)
	}

	file, err := os.Create(filepath.Base(filename))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	len, err = file.Write(filebytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len)
}

func SendFile(filename string, conn net.Conn) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	var filebytes []byte

	buffer := make([]byte, 1024)
	buffer = []byte(filename)
	_, err = conn.Write(buffer)
	if err != nil {
		fmt.Println(err)
	}

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
		filebytes = append(filebytes, buffer...)
	}

	conn.Write(filebytes)
}
