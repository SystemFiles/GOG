package command

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"sykesdev.ca/gog/lib"
)

func ExecPush() {
	var message string
	if len(flag.Args()) >= 2 {
		message = strings.Join(flag.Args()[1:], " ")
	}
	workingDir, GOGDir := lib.GetWorkspacePaths()

	feature, err := lib.NewFeatureFromFile()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to read feature from features file (%s). %v", workingDir + "/.gog/feature.json", err))
		os.Exit(1)
	}
	defer feature.Save()

	if message == "" {
		message = fmt.Sprintf("Feature Update (%d)", feature.TestCount)
		feature.UpdateTestCount()
	}

	if !lib.PathExists(GOGDir + "/feature.json") {
		lib.GetLogger().Error(fmt.Sprintf("Could not find valid GOG feature in the working directory. %s does not exist", GOGDir))
		os.Exit(1)
	}

	if err := lib.GitStageChanges(); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to stage existing changes. %v", err))
		os.Exit(1)
	}

	if !lib.GitHasUncommittedChanges() {
		lib.GetLogger().Warn(fmt.Sprintf("No un-committed changes were found for the current feature (%s).", feature.Jira))
	} else {
		if err := lib.GitCommitChanges(feature, message); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to commit changes to local project repo. %v", err))
			os.Exit(1)
		}
	}

	var pushArgs string
	if !lib.GitRemoteExists(feature) {
		pushArgs = fmt.Sprintf("--set-upstream origin %s", feature.Jira)
	}

	if err := lib.GitPushRemote(pushArgs); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to push changes to remote HEAD. %v", err))
		os.Exit(1)
	}

	lib.GetLogger().Info("Successfully pushed changes to remote feature!")
}