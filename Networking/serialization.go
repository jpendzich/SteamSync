package networking

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func SerializeFile(file Netfile, writer io.Writer) {
	SerializeString(file.Name, writer)
	SerializeInt(file.Length, writer)
	writer.Write(file.Actfile.Bytes())
	writer.Write(BoolToByte(file.Error))
}

func SerializeString(str Netstring, writer io.Writer) {
	SerializeInt(str.Length, writer)
	writer.Write([]byte(str.Actstr))
	writer.Write(BoolToByte(str.Error))
}

func SerializeInt(value uint64, writer io.Writer) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, value)
	writer.Write(buf)
}

func DeserializeFile(reader io.Reader) (Netfile, error) {
	var file Netfile
	var err error
	file.Name, err = DeserializeString(reader)
	file.Length = DeserializeInt(reader)
	file.Actfile = bytes.NewBuffer(nil)
	_, err = io.CopyN(file.Actfile, reader, int64(file.Length))
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1)
	reader.Read(buf)
	file.Error = ByteToBool(buf)
	if file.Error {
		err = errors.New("an error occurred on the other side of the connection")
	}
	return file, err
}

func DeserializeString(reader io.Reader) (Netstring, error) {
	var str Netstring
	var err error
	str.Length = DeserializeInt(reader)
	buf := make([]byte, str.Length)
	reader.Read(buf)
	str.Actstr = string(buf)
	buf = make([]byte, 1)
	reader.Read(buf)
	str.Error = ByteToBool(buf)
	if str.Error {
		err = errors.New("an error occurred on the other side of the connection")
	}
	return str, err
}

func DeserializeInt(reader io.Reader) uint64 {
	buf := make([]byte, 8)
	reader.Read(buf)
	return binary.BigEndian.Uint64(buf)
}
