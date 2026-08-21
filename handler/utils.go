package handler

import (
	"path/filepath"
	"os"
)

func GetExecutableDir() string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(exe)
}
