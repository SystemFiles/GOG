package command

import (
	"fmt"
	"os"
	"strings"

	"sykesdev.ca/gog/lib"
)

func ExecFinish(major, minor, patch bool) {
	_, GOGDir := lib.GetWorkspacePaths()

	if !major && !minor && !patch {
		lib.GetLogger().Error("You must specifiy exactly one of (-major, -minor, -patch) to create a new release")
		os.Exit(2)
	}

	feature, err := lib.NewFeatureFromFile()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to read feature from associated feature file. %v", err))
		os.Exit(1)
	}

	currentVersion, err := lib.GitOriginCurrentVersion()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to get current project version. %v", err))
		os.Exit(1)
	}

	updatedVersion := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(lib.BumpVersion(currentVersion, major, minor, patch))), "."), "[]")
	
	changelogEntry := lib.NewChangelogEntry(feature, updatedVersion, major)

	changelogLines, err := lib.CreateChangeLogLines(changelogEntry)
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to update the changelog. %v", err))
		os.Exit(1)
	}

	fmt.Println(changelogLines)

	if err := lib.WriteChangelogToFile(changelogLines); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to write changelog entry. %v", err))
		os.Exit(1)
	}

	if err := os.RemoveAll(GOGDir); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to remove GOG directory. %v", err))
		os.Exit(1)
	}

	if stderr, err := lib.GitStageChanges(); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to stage existing changes. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if stderr, err := lib.GitCommitChanges(feature, fmt.Sprintf("%s Final changes", feature.Jira)); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to commit changes to local project repo. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
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

	stderr, err := lib.GitCheckoutDefaultBranch()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to checkout default branch. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if stderr, err := lib.GitRebase(feature); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to rebase commits into new release. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if stderr, err := lib.GitCreateReleaseTags(updatedVersion, feature); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to create release tags. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if stderr, err := lib.GitPushRemote(""); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to push rebase to remote. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if stderr, err := lib.GitPushRemoteTagsOnly(); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to publish release tags to remote. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if stderr, err := feature.DeleteBranch(); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to delete existing feature branch for %s. %v", feature.Jira, err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	lib.GetLogger().Info(fmt.Sprintf("Successfully created new feature release for %s!", feature.Jira))
}