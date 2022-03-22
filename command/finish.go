package command

import (
	"fmt"
	"os"

	"sykesdev.ca/gog/lib"
	"sykesdev.ca/gog/lib/changelog"
	"sykesdev.ca/gog/lib/semver"
)

func bumpReleaseVersion(currentVersion semver.Semver, level semver.UpdateLevel) (semver.Semver) {

	switch level {
	case "MAJOR":
		return currentVersion.BumpMajor()
	case "MINOR":
		return currentVersion.BumpMinor()
	case "PATCH":
		return currentVersion.BumpPatch()
	default:
		return currentVersion
	}
}

func FinishUsage() {
	lib.GetLogger().Info("Usage: gog finish (-major || -minor || -patch)")
}

func ExecFinish(versionLevel semver.UpdateLevel) {
	workingDir, GOGDir := lib.WorkspacePaths()

	if !lib.GitIsValidRepo() {
		lib.GetLogger().Error(fmt.Sprintf("The current directory (%s) is not a valid git repository", workingDir))
		os.Exit(1)
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

	updatedVersion := bumpReleaseVersion(currentVersion, versionLevel)
	
	changelogEntry := changelog.NewChangelogEntry(feature, updatedVersion, versionLevel == "MAJOR")

	changelogLines, err := changelog.CreateChangeLogLines(changelogEntry)
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to update the changelog. %v", err))
		os.Exit(1)
	}

	if err := changelog.WriteChangelogToFile(changelogLines); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to write changelog entry. %v", err))
		os.Exit(1)
	}

	if err := os.RemoveAll(GOGDir); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to remove GOG directory. %v", err))
		os.Exit(1)
	}

	if stderr, err := feature.PushChanges(feature.Comment); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to push changes to remote repository. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	stderr, err := lib.GitCheckoutDefaultBranch()
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to checkout default branch. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if stderr, err := feature.Rebase(); err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Failed to rebase commits into new release. %v", err))
		lib.GetLogger().Error(stderr)
		os.Exit(1)
	}

	if stderr, err := feature.CreateReleaseTags(updatedVersion); err != nil {
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