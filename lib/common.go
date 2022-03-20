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

func CleanFeature(feature *Feature) error {
	_, GOGDir := GetWorkspacePaths()

	if _, err := GitCheckoutDefaultBranch(); err != nil {
		return err
	}

	cmd := exec.Command("git", "branch", "-D", feature.Jira)
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(GOGDir); err != nil {
		return err
	}

	return nil
}

func BumpVersion(currentVersion [3]int, major, minor, patch bool) ([3]int) {
	if major {
		currentVersion[0] += 1
		currentVersion[1] = 0
		currentVersion[2] = 0
		return currentVersion
	}

	if minor {
		currentVersion[1] += 1
		currentVersion[2] = 0
		return currentVersion
	}

	if patch {
		currentVersion[2] += 1
		return currentVersion
	}

	return [3]int{}
}