package gencli

import (
	"fmt"
	"strings"
)

const (
	// ExpectedParams is the number of expected plugin parameters
	ExpectedParams = 3
)

func parseParameters(params *string) (rootDir string, pbDir string, gapicDir string, err error) {
	if params == nil {
		err = fmt.Errorf("Missing required parameters. See usage")
		return
	}

	split := strings.Split(*params, ",")
	if len(split) != ExpectedParams {
		err = fmt.Errorf("Improper number of parameters. Got %d, require %d. See usage", len(split), ExpectedParams)
		return
	}

	for _, str := range split {
		sepNdx := strings.Index(str, ":")
		if sepNdx == -1 {
			err = fmt.Errorf("Unknown parameter: %s", str)
			return
		}

		switch str[:sepNdx] {
		case "grpc":
			pbDir = str[sepNdx+1:]
		case "gapic":
			gapicDir = str[sepNdx+1:]
		case "root":
			rootDir = str[sepNdx+1:]
		default:
			err = fmt.Errorf("Unknown parameter: %s", str)
		}
	}

	return
}

func strContains(a []string, s string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}

func extractMessageName(name string) string {
	return name[strings.LastIndex(name, ".")+1:]
}
