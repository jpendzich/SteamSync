package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func ReadString(conn net.Conn) string {
	writer := bytes.NewBufferString("")
	_, err := io.Copy(writer, conn)
	if err != nil {
		fmt.Println(err)
	}
	return string(writer.Bytes())
}

func WriteString(content string, conn net.Conn) {
	reader := bytes.NewBufferString(content)
	io.Copy(conn, reader)
}

func ReadFile(path string, conn net.Conn) {
	filename := ReadString(conn)
	path = filepath.Join(".", path)
	os.MkdirAll(path, os.ModePerm)
	file, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err1 := io.Copy(file, conn)
	if err1 != nil {
		fmt.Println(err)
	}
}

func SendFile(filename string, conn net.Conn) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	WriteString(filepath.Base(filename), conn)

	_, err1 := io.Copy(conn, file)
	if err1 != nil {
		fmt.Println(err)
	}
}
