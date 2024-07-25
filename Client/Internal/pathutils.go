package internal

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func BuildWindowsPath(path string) string {
	fullpath := ""
	regex, err := regexp.Compile(`%\w*%`)
	if err != nil {
		log.Println(err)
	}
	env := string(regex.Find([]byte(path)))
	envval := strings.ReplaceAll(env, "%", "")
	envval = os.Getenv(envval)

	fullpath = strings.Replace(path, env, envval, 1)

	return fullpath
}

func BuildSteamDeckPath(linpath string, winpath string) string {
	windowsEnvVars := make(map[string]string)
	windowsEnvVars["USERPROFILE"] = "drive_c/users/steamuser/"
	windowsEnvVars["APPDATA"] = "drive_c/users/steamuser/AppData/Roaming"
	windowsEnvVars["LOCALAPPDATA"] = "drive_c/users/steamuser/AppData/Local"

	fullpath := ""

	fullpath = strings.Replace(linpath, "<SteamLibrary-folder>", "/home/deck/.steam", 1)

	winpath = strings.ReplaceAll(winpath, "\\", "/")

	regex, err := regexp.Compile(`%\w*%`)
	if err != nil {
		log.Println(err)
	}
	env := string(regex.Find([]byte(winpath)))
	envval := strings.ReplaceAll(env, "%", "")
	log.Println(envval)
	envval = windowsEnvVars[strings.ToUpper(envval)]
	fullpath = fullpath + strings.Replace(winpath, env, envval, 1)

	return fullpath
}
