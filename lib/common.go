package lib

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

func GetWorkspacePaths() (string, string) {
	workingDir, err := os.Getwd()
	if err != nil {
		GetLogger().Fatal("Failed to get working directory from path")
	}
	GOGDir := workingDir + "/.gog"

	return workingDir, GOGDir
}