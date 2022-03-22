package command

import (
	"fmt"
	"os"

	"sykesdev.ca/gog/lib"
)

func PushUsage() {
	lib.GetLogger().Info("Usage: gog push [message ...]")
}

func ExecPush(message string) {
	workingDir, _ := lib.WorkspacePaths()

	if !lib.GitIsValidRepo() {
		lib.GetLogger().Error(fmt.Sprintf("The current directory (%s) is not a valid git repository", workingDir))
		os.Exit(1)
	}

	feature, err := lib.NewFeatureFromFile()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to read feature from features file (%s). %v", workingDir + "/.gog/feature.json", err))
		os.Exit(1)
	}
	defer feature.Save()
	
	if message == "" {
		message = fmt.Sprintf("Test Build (%d)", feature.TestCount)
		feature.UpdateTestCount()
	}

	if stderr, err := lib.GitPublishChanges(feature, message); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to publish changes to remote repository. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	lib.GetLogger().Info("Successfully pushed changes to remote feature!")
}