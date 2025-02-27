package internal

import (
	"os"
	"path/filepath"
)

func GetSteamInstallPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	installPath := filepath.Join(home, ".local/share/Steam/")
	return installPath, nil
}
