package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"sykesdev.ca/gog/changelog"
	"sykesdev.ca/gog/common"
	"sykesdev.ca/gog/config"
	"sykesdev.ca/gog/git"
	"sykesdev.ca/gog/logging"
	"sykesdev.ca/gog/models"
	"sykesdev.ca/gog/prompt"
	"sykesdev.ca/gog/semver"
)

type FinishAction string

func bumpReleaseVersion(currentVersion semver.Semver, action FinishAction) (semver.Semver) {
	switch action {
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

type FinishCommand struct {
	fs *flag.FlagSet

	name string
	alias string
	action FinishAction

	major bool
	minor bool
	patch bool
}

func NewFinishCommand() *FinishCommand {
	fc := &FinishCommand{
		name: "finish",
		alias: "fin",
		fs: flag.NewFlagSet("finish", flag.ContinueOnError),
	}

	fc.fs.BoolVar(&fc.major, "major", false, "specifies that in this freature you make incompatible API changes (breaking changes)")
	fc.fs.BoolVar(&fc.minor, "minor", false, "specifies that in this feature you add functionality in a backwards compatible manner (non-breaking)")
	fc.fs.BoolVar(&fc.patch, "patch", false, "specifies that in this feature you make backwards compatible bug fixes small backwards compatible updates")

	fc.fs.Usage = fc.Help

	return fc
}

func (fc *FinishCommand) Help() {
	fmt.Printf(
`Usage: %s (%s | %s) (-major | -minor | -patch) [-h] [-help]

-------====== Finish Arguments ======-------

`, os.Args[0], fc.name, fc.alias)

	fc.fs.PrintDefaults()

	fmt.Println("\n-------================================-------")
}

func (fc *FinishCommand) Init(args []string) error {
	err := fc.fs.Parse(args)

	if fc.major {
		fc.action = "MAJOR"
	} else if fc.minor {
		fc.action = "MINOR"
	} else if fc.patch {
		fc.action = "PATCH"
	}

	if fc.action == "" {
		return errors.New("failed to specify major, minor or patch for this feature upgrade (re-run wiht -h for full usage details)")
	}

	return err
}

func (fc *FinishCommand) Run() error {
	workingDir, GOGDir := common.WorkspacePaths()

	if !git.IsValidRepo() {
		return fmt.Errorf("the current directory (%s) is not a valid git repository", workingDir)
	}

	if !common.PathExists(GOGDir + "/feature.json") {
		return errors.New("feature file not found ... there may not be a GOG feature on this branch")
	}

	feature, err := models.NewFeatureFromFile()
	if err != nil {
		return fmt.Errorf("failed to read feature from associated feature file. %v", err)
	}

	if feature.CustomVersionPrefix != config.AppConfig().TagPrefix() && feature.CustomVersionPrefix != "" {
		config.AppConfig().SetTagPrefix(feature.CustomVersionPrefix)
	}

	existingPrefix, err := git.ProjectExistingVersionPrefix()
	if err != nil {
		return fmt.Errorf("failed to get projects existing version prefix. %v", err)
	}
	if existingPrefix != config.AppConfig().TagPrefix() {
		logging.GetLogger().Warn(fmt.Sprintf("feature version prefix specified does not match existing prefix for this git project ('%s' != '%s')", config.AppConfig().TagPrefix(), existingPrefix))
		if c := prompt.String("continue with feature release (Y/n)? "); strings.ToUpper(c) != "Y" {
			logging.GetLogger().Info("safely exiting feature release")
			return nil
		}
		logging.GetLogger().Info("continuing with feature release against warning")
	}

	currentVersion, err := git.OriginCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current project version. %v", err)
	}

	updatedVersion := bumpReleaseVersion(currentVersion, fc.action)

	changelogEntry := changelog.NewChangelogEntry(feature, updatedVersion, fc.action == "MAJOR" || fc.action == "MINOR")

	changelogLines, err := changelog.CreateChangeLogLines(changelogEntry)
	if err != nil {
		return fmt.Errorf("failed to update the changelog. %v", err)
	}

	if err := changelog.WriteChangelogToFile(changelogLines); err != nil {
		return fmt.Errorf("failed to write changelog entry. %v", err)
	}

	if err := os.RemoveAll(GOGDir); err != nil {
		return fmt.Errorf("failed to remove GOG directory. %v", err)
	}

	if stderr, err := feature.PushChanges(feature.Comment); err != nil {
		return fmt.Errorf("failed to push changes to remote repository. %v \n%s", err, stderr)
	}

	stderr, err := git.CheckoutDefaultBranch()
	if err != nil {
		return fmt.Errorf("failed to checkout default branch. %v \n%s", err, stderr)
	}

	if stderr, err := feature.Rebase(); err != nil {
		return fmt.Errorf("failed to rebase commits into new release. %v \n%s", err, stderr)
	}

	if stderr, err := feature.CreateReleaseTags(updatedVersion); err != nil {
		return fmt.Errorf("failed to create release tags. %v \n%s", err, stderr)
	}

	if stderr, err := git.PushRemote(""); err != nil {
		return fmt.Errorf("failed to push rebase to remote. %v \n%s", err, stderr)
	}

	if stderr, err := git.PushRemoteTagsOnly(); err != nil {
		return fmt.Errorf("failed to publish release tags to remote. %v \n%s", err, stderr)
	}

	if stderr, err := feature.DeleteBranch(); err != nil {
		return fmt.Errorf("failed to delete existing feature branch for %s. %v \n%s", feature.Jira, err, stderr)
	}

	logging.GetLogger().Info(fmt.Sprintf("Successfully created new feature release for %s!", feature.Jira))

	return nil
}

func (fc *FinishCommand) Name() string {
	return fc.name
}

func (fc *FinishCommand) Alias() string {
	return fc.alias
}