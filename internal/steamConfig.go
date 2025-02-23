package internal

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/andygrunwald/vdf"
	"golang.org/x/sys/windows/registry"
)

const (
	steamApps         = "steamapps"
	common            = "common"
	appManifestPrefix = "appmanifest_"
	appManifestSuffix = ".acf"
)

func GetSteamInstallPath() (string, error) {
	sys := runtime.GOOS

	installPath := ""
	switch sys {
	case "windows":
		key, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Valve\\Steam", registry.QUERY_VALUE)
		if err != nil {
			return "", err
		}
		defer key.Close()

		installPath, _, err = key.GetStringValue("SteamPath")
		if err != nil {
			return "", err
		}
	case "linux":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		installPath = filepath.Join(home, ".local/share/Steam/")
	}
	return installPath, nil
}

func GetGameInstallPath(appId int) (string, error) {
	steamInstallPath, err := GetSteamInstallPath()
	if err != nil {
		return "", err
	}
	libraryfolders, err := os.Open(filepath.Join(steamInstallPath, steamApps, "libraryfolders.vdf"))
	parser := vdf.NewParser(libraryfolders)
	content, err := parser.Parse()
	if err != nil {
		return "", err
	}
	libraries := content["libraryfolders"].(map[string]interface{})
	libraryPath, err := findLibrary(appId, libraries)
	if err != nil {
		return "", err
	}

	appManifestName := appManifestPrefix + strconv.Itoa(appId) + appManifestSuffix
	appManifest, err := os.Open(filepath.Join(libraryPath, steamApps, appManifestName))
	parser = vdf.NewParser(appManifest)
	content, err = parser.Parse()
	if err != nil {
		return "", err
	}
	installDir, err := findInstallDir(content["AppState"].(map[string]interface{}))
	if err != nil {
		return "", err
	}

	installPath := filepath.Join(libraryPath, steamApps, common, installDir)

	return installPath, nil
}

func findLibrary(appId int, libraries map[string]interface{}) (string, error) {
	for _, library := range libraries {
		apps := library.(map[string]interface{})["apps"].(map[string]interface{})
		for app := range apps {
			appAsInt, err := strconv.Atoi(app)
			if err != nil {
				return "", err
			}
			if appAsInt == appId {
				return library.(map[string]interface{})["path"].(string), nil
			}
		}
	}
	return "", nil
}

func findInstallDir(appState map[string]interface{}) (string, error) {
	installdir, ok := appState["installdir"].(string)
	if ok {
		return installdir, nil
	}
	return "", errors.New("Key installdir not found in file")
}
