package utils

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

func TestDeserializeInt(test *testing.T) {
	input := make([]byte, 8)
	binary.BigEndian.PutUint64(input, 23)
	input = append(input, input...)
	output := DeserializeInt(&input)
	if output.payload != 23 {
		test.Fatalf("Integer deserialization failed")
	}
}

func TestDeserializeString(test *testing.T) {
	input := []byte("Hello, this is a String")
	inputlength := make([]byte, 8)
	binary.BigEndian.PutUint64(inputlength, uint64(len(input)))
	input = append(inputlength, input...)
	input = append(input, input...)
	output := DeserializeString(&input)
	if output.payload != "Hello, this is a String" {
		test.Fatalf("String deserialization failed: %s : %d", output.payload, output.length)
	}
}

func TestDeserializeFile(test *testing.T) {
	filename := "test.txt"
	file, _ := os.Open(filename)
	defer file.Close()
	stat, _ := file.Stat()
	filebuf := make([]byte, stat.Size())
	file.Read(filebuf)
	buf := filebuf
	buflen := make([]byte, 8)
	binary.BigEndian.PutUint64(buflen, uint64(stat.Size()))
	buf = append(buflen, buf...)
	buf = append([]byte(filename), buf...)
	binary.BigEndian.PutUint64(buflen, uint64(len(filename)))
	buf = append(buflen, buf...)
	output := DeserializeFile(&buf)
	if !bytes.Equal(output.payload, filebuf) {
		test.Fatalf("File deserialization failed: %s\n %s", output.payload, filebuf)
	}
}

func TestSerializeInt(test *testing.T) {
	output := SerializeInt(PInt{23})
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 23)
	if !bytes.Equal(output, buf) {
		test.Fatalf("Int serialization failed %d and %d", output, buf)
	}
}

func TestSerializeString(test *testing.T) {
	output := SerializeString(PString{PInt{23}, "Hello, this is a String"})
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 23)
	buf = append(buf, []byte("Hello, this is a String")...)
	if !bytes.Equal(output, buf) {
		test.Fatalf("String serialization failed")
	}
}

func TestSerializeFile(test *testing.T) {
	filename := "test.txt"
	file, _ := os.Open(filename)
	defer file.Close()
	stat, _ := file.Stat()
	filebuf := make([]byte, stat.Size())
	file.Read(filebuf)
	var packet PFile
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(len(filename)))
	buf = append(buf, []byte(filename)...)
	buflen := make([]byte, 8)
	binary.BigEndian.PutUint64(buflen, uint64(stat.Size()))
	buf = append(buf, buflen...)
	buf = append(buf, filebuf...)
	packet.lengthpacket.payload = uint64(len(buf))
	packet.name.length.payload = uint64(len(filename))
	packet.name.payload = filename
	packet.lengthpayload.payload = uint64(stat.Size())
	packet.payload = filebuf
	output := SerializeFile(packet)
	if !bytes.Equal(output, buf) {
		test.Fatalf("File serialization failed %s \n %s", output, buf)
	}
}
