package networking

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func WriteFile(file Netfile, writer io.Writer) {
	WriteError(false, writer)
	WriteString(file.Name, writer)
	WriteInt(file.Length, writer)
	writer.Write(file.Actfile.Bytes())
}

func WriteString(str Netstring, writer io.Writer) {
	WriteError(false, writer)
	WriteInt(str.Length, writer)
	writer.Write([]byte(str.Actstr))
}

func WriteInt(value uint64, writer io.Writer) {
	WriteError(false, writer)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, value)
	writer.Write(buf)
}

func WriteError(hasError bool, writer io.Writer) {
	writer.Write(BoolToByte(hasError))
}

func ReadFile(reader io.Reader) (Netfile, error) {
	var file Netfile
	var err error
	if ReadError(reader) {
		return file, errors.New("an error occured on the other side of the connection")
	}
	file.Name, err = ReadString(reader)
	if err != nil {
		return file, err
	}
	file.Length, err = ReadInt(reader)
	if err != nil {
		return file, err
	}
	file.Actfile = bytes.NewBuffer(nil)
	_, err = io.CopyN(file.Actfile, reader, int64(file.Length))
	if err != nil {
		return file, err
	}
	return file, err
}

func ReadString(reader io.Reader) (Netstring, error) {
	var str Netstring
	var err error
	if ReadError(reader) {
		return str, errors.New("an error occured on the other side of the connection")
	}
	str.Length, err = ReadInt(reader)
	if err != nil {
		return str, err
	}
	buf := make([]byte, str.Length)
	reader.Read(buf)
	str.Actstr = string(buf)
	return str, nil
}

func ReadInt(reader io.Reader) (uint64, error) {
	if ReadError(reader) {
		return 0, errors.New("an error occured on the other side of the connection")
	}
	buf := make([]byte, 8)
	reader.Read(buf)
	return binary.BigEndian.Uint64(buf), nil
}

func ReadError(reader io.Reader) bool {
	buf := make([]byte, 1)
	reader.Read(buf)
	return ByteToBool(buf)
}
