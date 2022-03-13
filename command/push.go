package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"sykesdev.ca/gog/lib"
)

func ReadFeature(cwd string) (*lib.Feature, error) {
	featureBytes, err := os.ReadFile(cwd + "/.gog/feature.json")
	if err != nil {
		return nil, err
	}

	var feature *lib.Feature
	err = json.Unmarshal(featureBytes, &feature)
	if err != nil {
		return nil, err
	}

	return feature, nil
}

func GitStageChanges() error {
	cmd := exec.Command("git", "add", "-A")
	_, err := cmd.Output()
	
	return err
}

func GitCommitChanges(feature *lib.Feature, commitMessage string) error {
	formattedMessage := fmt.Sprintf("%s %s", feature.Jira, commitMessage)
	cmd := exec.Command("git", "commit", "-m", formattedMessage)
	_, err := cmd.Output()
	
	return err
}

func GitRemoteExists(feature *lib.Feature) bool {
	remoteExistsCommand := fmt.Sprintf("git ls-remote --heads --exit-code | egrep %s", feature.Jira)
	cmd := exec.Command("bash", "-c", remoteExistsCommand)
	_, err := cmd.Output()

	return err == nil
}

func GitPushRemote(pushArgs string) error {
	pushCommand := fmt.Sprintf("git push %s", pushArgs)
	cmd := exec.Command("bash", "-c", pushCommand)
	_, err := cmd.Output()

	return err
}

func ExecPush() {
	var message string
	if len(flag.Args()) >= 2 {
		message = strings.Join(flag.Args()[1:], " ")
	}
	workingDir, GOGDir := lib.GetWorkspacePaths()

	feature, err := ReadFeature(workingDir)
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to read feature from features file (%s). %v", workingDir + "/.gog/feature.json", err))
		os.Exit(1)
	}
	defer feature.Save()

	if message == "" {
		message = fmt.Sprintf("Feature Update (%d)", feature.TestCount)
		feature.UpdateTestCount()
	}

	if !PathExists(GOGDir + "/feature.json") {
		lib.GetLogger().Error(fmt.Sprintf("Could not find valid GOG feature in the working directory. %s does not exist", GOGDir))
		os.Exit(1)
	}

	if err := GitStageChanges(); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to stage existing changes. %v", err))
		os.Exit(1)
	}

	if !lib.GitHasUncommittedChanges() {
		lib.GetLogger().Warn(fmt.Sprintf("No un-committed changes were found for the current feature (%s).", feature.Jira))
	} else {
		if err := GitCommitChanges(feature, message); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to commit changes to local project repo. %v", err))
			os.Exit(1)
		}
	}

	var pushArgs string
	if !GitRemoteExists(feature) {
		pushArgs = fmt.Sprintf("--set-upstream origin %s", feature.Jira)
	}

	if err := GitPushRemote(pushArgs); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to push changes to remote HEAD. %v", err))
		os.Exit(1)
	}

	lib.GetLogger().Info("Successfully pushed changes to remote feature!")
}