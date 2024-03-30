package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func ReadString(conn net.Conn) string {
	buflen := make([]byte, 8)
	conn.Read(buflen)
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	io.CopyN(writer, conn, int64(binary.BigEndian.Uint64(buflen)))
	buflen = append(buflen, buf.Bytes()...)
	packet := DeserializeString(&buflen)
	return packet.payload
}

func WriteString(content string, conn net.Conn) {
	var packet PString
	packet.payload = content
	packet.length.payload = uint64(len((content)))
	conn.Write(SerializeString(packet))
}

func ReadFile(path string, conn net.Conn) {
	filename := ReadString(conn)

	fileLength := make([]byte, 8)
	_, err := conn.Read(fileLength)
	if err != nil {
		fmt.Println(err)
	}

	length := binary.LittleEndian.Uint64(fileLength)
	fmt.Println(length)
	fileBytes := make([]byte, length)
	_, err1 := conn.Read(fileBytes)
	if err1 != nil {
		fmt.Println(err)
	}
	path = filepath.Join(".", path)
	os.MkdirAll(path, os.ModePerm)
	file, err := os.Create(filepath.Join(path, filename))
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

	var packet PFile
	packet.name.length.payload = uint64(len(filename))
	packet.name.payload = filename
	packet.length.payload = uint64(fileInfo.Size())
	buf := make([]byte, fileInfo.Size())
	file.Read(buf)
	packet.payload = buf
	_, err = conn.Write(SerializeFile(packet))
	if err != nil {
		fmt.Println(err)
	}
}
