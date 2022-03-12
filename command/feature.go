package command

import (
	"fmt"
	"os"

	"sykesdev.ca/gog/lib"
)

func usage() {
	lib.GetLogger().Info("Usage: gog feature <jira_name> [comment]")
}

func ExecFeature(opts []string) {
	if len(opts) < 1 {
		usage()
		os.Exit(0)
	}

	jira := opts[0]
	comment := opts[1]
	if comment == "" {
		comment = "Feature Branch"
	}

	feature, err := lib.NewFeature(jira, comment)
	if err != nil {
		lib.GetLogger().Error("Failed to create feature")
		lib.GetLogger().Error(fmt.Sprintf("Reason: %v", err))
		os.Exit(1)
	}

	lib.GetLogger().Info(fmt.Sprintf("Successfully created feature %s!", jira))
}