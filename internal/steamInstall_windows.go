package internal

import "golang.org/x/sys/windows/registry"

func GetSteamInstallPath() (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Valve\\Steam", registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	installPath, _, err := key.GetStringValue("SteamPath")
	if err != nil {
		return "", err
	}
	return installPath, nil
}
