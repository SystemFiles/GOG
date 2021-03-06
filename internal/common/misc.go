package common

import (
	"os"
	"os/exec"
)

func GitProjectRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return CleanStdoutSingleline(stdout), nil
}

func StringInSlice(slice []string, value string) (bool) {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func GOGPath() (string) {
	projectRoot, err := GitProjectRoot()
	if err != nil {
		panic("cannot determine GOG configuration path since we cannot find the root of this git repo")
	}
	GOGDir := projectRoot + "/.gog"

	return GOGDir
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}