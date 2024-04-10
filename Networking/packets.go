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
	Error   bool
}

type Netstring struct {
	Length uint64
	Actstr string
	Error  bool
}

func BuildNetstring(str string) Netstring {
	var netstr Netstring
	netstr.Length = uint64(len(str))
	netstr.Actstr = str
	return netstr
}

func BuildNetfile(file *os.File) Netfile {
	var netfile Netfile
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

func BoolToByte(b bool) byte {
	if b {
		return 1
	} else {
		return 0
	}
}

func ByteToBool(b byte) bool {
	return b == 1
}
