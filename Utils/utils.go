package utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func ReadString(conn net.Conn) string {
	lengthBytes := make([]byte, 8)
	_, err := conn.Read(lengthBytes)
	if err != nil {
		fmt.Println(err)
	}
	length := binary.LittleEndian.Uint64(lengthBytes)

	contentBytes := make([]byte, length)
	_, err1 := conn.Read(contentBytes)
	if err1 != nil {
		fmt.Println(err1)
	}

	return string(contentBytes)
}

func WriteString(content string, conn net.Conn) {
	contentBytes := []byte(content)
	lengthBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lengthBytes, uint64(len(contentBytes)))

	_, err := conn.Write(lengthBytes)
	if err != nil {
		fmt.Println(err)
	}
	_, err1 := conn.Write(contentBytes)
	if err1 != nil {
		fmt.Println(err)
	}
}

func ReadFile(conn net.Conn) {
	fileLength := make([]byte, 8)
	_, err := conn.Read(fileLength)
	if err != nil {
		fmt.Println(err)
	}

	length := binary.LittleEndian.Uint64(fileLength)
	fileBytes := make([]byte, length)
	_, err1 := conn.Read(fileBytes)
	if err1 != nil {
		fmt.Println(err)
	}

	file, err := os.Create("TEST.md")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = file.Write(fileBytes)
	if err != nil {
		fmt.Println(err)
	}
}

func SendFile(filename string, conn net.Conn) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	fileLength := make([]byte, 8)
	binary.LittleEndian.PutUint64(fileLength, uint64(fileInfo.Size()))
	_, err = conn.Write(fileLength)
	if err != nil {
		fmt.Println(err)
	}

	filebytes := make([]byte, fileInfo.Size())
	_, err1 := file.Read(filebytes)
	if err1 != nil {
		fmt.Println(err)
	}

	_, err = conn.Write(filebytes)
	if err != nil {
		fmt.Println(err)
	}

	conn.Write(filebytes)
}
