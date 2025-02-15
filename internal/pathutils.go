package internal

import (
	"os"
	"regexp"
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
	}
)

func ConvertToRealPath(wikiPath string) (string, error) {
	regex, err := regexp.Compile("{{p|.*?}}")
	if err != nil {
		return "", err
	}

	newPath := replaceAllStringSubmatchFunc(regex, wikiPath, func(groups []string) string {
		replacedString := ""
		if windowsEnvVars[groups[1]] {
			replacedString = os.Getenv(groups[1])
		} else if linuxEnvVars[groups[1]] {
			switch groups[1] {
			case "home":
				replacedString, err = os.UserHomeDir()
				if err != nil {
					// TODO: will be an invalid path
					return groups[1]
				}
			case "linuxhome":
				replacedString, err = os.UserHomeDir()
				if err != nil {
					// TODO: will be an invalid path
					return groups[1]
				}

			}
		} else if allEnvVars[groups[1]] {
			switch groups[1] {
			case "steam":
				// TODO: Steam install dir
			case "uid":
				// TODO: Get the users steam id
			}
		}
		return replacedString
	})

	return newPath, nil
}

// copied from https://gist.github.com/elliotchance/d419395aa776d632d897 for convenience
func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}
