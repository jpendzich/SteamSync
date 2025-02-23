package internal

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The paths in the Pcgamingwiki save game locations are often made up of environment variables or
// other kind of variables. These are mostly identifiably by the {{p|*}} syntax, with * being the
// name of the variable. The maps below aim to group these into different platforms with some
// being independant of platform like the "steam" and "uid" with are the steam install directory
// and the steam user id respectively

var (
	windowsEnvVars = map[string]bool{
		"userprofile":  true,
		"appdata":      true,
		"localappdata": true,
		"programdata":  true,
		"programfiles": true,
		"public":       true,
	}

	linuxEnvVars = map[string]bool{
		"home":      true,
		"linuxhome": true,
	}

	allEnvVars = map[string]bool{
		"steam": true,
		"uid":   true,
		"game":  true,
	}
)

func ConvertToRealPath(appId int, wikiPath string) (string, error) {
	// match {{p|*}} or {{P|*}} and return * as capture group
	regex, err := regexp.Compile("{{[p|P]\\|(.*?)}}")
	if err != nil {
		return "", err
	}

	newPath, err := replaceAllStringSubmatchFunc(regex, wikiPath, func(groups []string) (string, error) {
		replacedString := ""
		envVarAndPath := strings.SplitN(groups[1], "\\", 2)
		if len(envVarAndPath) < 1 {
			envVarAndPath = strings.SplitN(groups[1], "/", 2)
		}
		envVar := envVarAndPath[0]
		restPath := ""
		if len(envVarAndPath) > 1 {
			restPath = envVarAndPath[1]
		}

		if windowsEnvVars[envVar] {
			replacedString = filepath.Join(os.Getenv(envVar), restPath)
		} else if linuxEnvVars[envVar] {
			switch envVar {
			case "home":
				replacedString, err = os.UserHomeDir()
				if err != nil {
					return "", err
				}
				replacedString = filepath.Join(replacedString, restPath)
			case "linuxhome":
				replacedString, err = os.UserHomeDir()
				if err != nil {
					return "", err
				}
				replacedString = filepath.Join(replacedString, restPath)
			}
		} else if allEnvVars[envVar] {
			switch envVar {
			case "steam":
				replacedString, err = GetSteamInstallPath()
				if err != nil {
					return "", err
				}
				replacedString = filepath.Join(replacedString, restPath)
			case "uid":
				// TODO: Get the users steam id
			case "game":
				replacedString, err = GetGameInstallPath(appId)
				if err != nil {
					return "", err
				}
				replacedString = filepath.Join(replacedString, restPath)
			}
		}
		return replacedString, nil
	})
	if err != nil {
		return "", err
	}

	return newPath, nil
}

// copied from https://gist.github.com/elliotchance/d419395aa776d632d897 for convenience
func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) (string, error)) (string, error) {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}
		replaced, err := repl(groups)
		if err != nil {
			return "", err
		}
		result += str[lastIndex:v[0]] + replaced
		lastIndex = v[1]
	}

	return result + str[lastIndex:], nil
}
