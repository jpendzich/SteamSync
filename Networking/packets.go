package networking

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

type Netfile struct {
	Name    Netstring
	Length  uint64
	Actfile *bytes.Buffer
}

type Netstring struct {
	Length uint64
	Actstr string
}

func BuildNetstring(str string) Netstring {
	var netstr Netstring
	netstr.Length = uint64(len(str))
	netstr.Actstr = str
	return netstr
}

func BuildNetfile(file *os.File) Netfile {
	var netfile Netfile
	if file == nil {
		netfile.Name = BuildNetstring("")
		netfile.Length = 0
		netfile.Actfile = bytes.NewBuffer(nil)
		return netfile
	}

	netfile.Name = BuildNetstring(filepath.Base(file.Name()))
	filestat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	netfile.Length = uint64(filestat.Size())
	netfile.Actfile = bytes.NewBuffer(nil)
	_, err = io.Copy(netfile.Actfile, file)
	if err != nil {
		panic(err)
	}
	return netfile
}

func BoolToByte(b bool) []byte {
	buf := make([]byte, 1)
	if b {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return buf
}

func ByteToBool(b []byte) bool {
	return b[0] == 1
}
