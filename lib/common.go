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

func WorkspacePaths() (string, string) {
	workingDir, err := os.Getwd()
	if err != nil {
		GetLogger().Fatal("Failed to get working directory from path")
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

func CleanFeature(feature *Feature) error {
	_, GOGDir := WorkspacePaths()

	if _, err := GitCheckoutDefaultBranch(); err != nil {
		return err
	}

	if _, err := feature.DeleteBranch(); err != nil {
		return err
	}

	if err := os.RemoveAll(GOGDir); err != nil {
		return err
	}

	return nil
}