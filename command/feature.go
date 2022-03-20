package command

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"sykesdev.ca/gog/lib"
)

func featureUsage() {
	lib.GetLogger().Info("Usage: gog feature <jira_name> <comment> [-from-feature]")
}

func ExecFeature(fromFeature bool) {
	if len(flag.Args()) < 3 {
		lib.GetLogger().Error("Invalid usage of gogfeature ...")
		featureUsage()
		os.Exit(2)
	}

	jira := flag.Arg(1)
	comment := strings.Join(flag.Args()[2:], " ")

	if comment == "" {
		comment = "Feature Branch"
	}

	validJiraFormat, err := regexp.Match(`^[A-Z].*\-[0-9].*$`, []byte(jira))
	if err != nil {
		lib.GetLogger().Fatal(fmt.Sprintf("Failed to parse regular expression for Jira format. %v", err))
	}

	if !validJiraFormat {
		lib.GetLogger().Error("Invalid Jira format ... example of a valid format includes JIRA-0023")
		os.Exit(1)
	}

	_, GOGDir := lib.GetWorkspacePaths()

	feature, err := lib.NewFeature(jira, comment)
	if err != nil {
		lib.GetLogger().Error("Failed to create feature")
		lib.GetLogger().Error(fmt.Sprintf("Reason: %v", err))
		os.Exit(1)
	}

	if !lib.GitIsValidRepo() {
		lib.GetLogger().Error("The current directory does not contain a valid git repo ...")
		os.Exit(1)
	}

	initial_branch, err := lib.GitGetCurrentBranch()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to get current branch of git repository. %v", err))
		os.Exit(1)
	}

	if feature.BranchExists() {
		lib.GetLogger().Error(fmt.Sprintf("There is already a branch in this repo named %s", feature.Jira))
		os.Exit(1)
	}

	if lib.GitHasUnstagedCommits() {
		lib.GetLogger().Error(fmt.Sprintf("There is unstaged commits on your current branch (%s). For your safety, please stage or discard the changes to continue.", initial_branch))
		os.Exit(1)
	}

	if lib.GitHasUncommittedChanges() {
		lib.GetLogger().Error(fmt.Sprintf("There are staged commits on your current branch (%s) which have not been committed. %v", initial_branch, err))
		os.Exit(1)
	}

	if !fromFeature {
		if stderr, err := lib.GitCheckoutDefaultBranch(); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to checkout default branch for repo. %v", err))
			lib.GetLogger().Error(stderr)
			os.Exit(1)
		}
	}

	if lib.PathExists(GOGDir) {
		lib.GetLogger().Error(fmt.Sprintf("%s already exists ... there could already be a feature here. Please fix this and try again.", GOGDir))
		os.Exit(1)
	}

	if stderr, err := feature.CreateBranch(true); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to create or checkout new feature branch, %s. %v", feature.Jira, err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if err := feature.Save(); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to create feature tracking file (%v) ... will exit cleanly", err))
		if err := lib.CleanFeature(feature); err != nil {
			lib.GetLogger().Error(fmt.Sprintf("Failed to exit cleanly ... %v", err))
		}
		os.Exit(1)
	}

	lib.GetLogger().Info(fmt.Sprintf("Successfully created feature %s!", jira))
}