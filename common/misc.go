package common

import (
	"os"
	"os/exec"
	"strings"
)

func gitLocalRepositoryRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(stdout), nil
}

func StringInSlice(slice []string, value string) (bool) {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func WorkspacePaths() (string, string) {
	workingDir, err := os.Getwd()
	if err != nil {
		panic("could not get current working directory")
	}
	
	repoRoot, err := gitLocalRepositoryRoot()
	if err != nil {
		panic("cannot determine GOG configuration path since we cannot find the root of this git repo")
	}
	GOGDir := strings.ReplaceAll(string(repoRoot), "\n", "") + "/.gog"

	return workingDir, GOGDir
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}