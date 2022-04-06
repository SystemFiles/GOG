package common

import (
	"os"
)

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
	GOGDir := workingDir + "/.gog"

	return workingDir, GOGDir
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}