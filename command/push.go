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
	workingDir, _ := lib.GetWorkspacePaths()

	feature, err := lib.NewFeatureFromFile()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to read feature from features file (%s). %v", workingDir + "/.gog/feature.json", err))
		os.Exit(1)
	}
	defer feature.Save()

	if message == "" {
		message = fmt.Sprintf("%s Test Build (%d)", feature.Jira, feature.TestCount)
		feature.UpdateTestCount()
	}

	if stderr, err := lib.GitStageChanges(); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to stage existing changes. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if !lib.GitHasUncommittedChanges() {
		lib.GetLogger().Warn(fmt.Sprintf("No un-committed changes were found for the current feature (%s).", feature.Jira))
	} else {
		if stderr, err := lib.GitCommitChanges(feature, message); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to commit changes to local project repo. %v", err))
			lib.GetLogger().Error(stderr)
			os.Exit(1)
		}
	}

	var pushArgs string
	if !feature.RemoteExists() {
		pushArgs = fmt.Sprintf("--set-upstream origin %s", feature.Jira)
	} else {
		// only pull changes if a remote exists
		if stderr, err := lib.GitPullChanges(); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to pull changes from remote before push. %v", err))
			lib.GetLogger().Error(stderr)
			os.Exit(1)
		}
	}

	if stderr, err := lib.GitPushRemote(pushArgs); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to push changes to remote HEAD. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	lib.GetLogger().Info("Successfully pushed changes to remote feature!")
}