package lib

import (
	"os"
	"os/exec"
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

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CleanFeature(cwd string, feature *Feature) error {
	if _, err := GitCheckoutDefaultBranch(); err != nil {
		return err
	}

	cmd := exec.Command("git", "branch", "-D", feature.Jira)
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	os.RemoveAll(cwd + "/.gog/")

	return nil
}