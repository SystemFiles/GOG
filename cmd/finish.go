package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"sykesdev.ca/gog/config"
	"sykesdev.ca/gog/internal/changelog"
	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/git"
	"sykesdev.ca/gog/internal/logging"
	"sykesdev.ca/gog/internal/models"
	"sykesdev.ca/gog/internal/prompt"
	"sykesdev.ca/gog/internal/semver"
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

	noChangelog bool
	noTag bool
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
	fc.fs.BoolVar(&fc.noChangelog, "no-changelog", false, "if this flag is set, no changelog creation or updates shall be performed when finishing this feature release")
	fc.fs.BoolVar(&fc.noTag, "no-tag", false, "if this flag is set, no version tagging shall be applied to this finished feature release")	

	fc.fs.Usage = fc.Help

	return fc
}

func (fc *FinishCommand) Help() {
	fmt.Printf(
`Usage: %s (%s | %s) (-major | -minor | -patch) [ additional_options... ] [-h] [-help]

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
	GOGDir := common.GOGPath()

	if !common.PathExists(GOGDir + "/feature.json") {
		return errors.New("feature file not found ... there may not be a GOG feature on this branch")
	}

	feature, err := models.NewFeatureFromFile()
	if err != nil {
		return fmt.Errorf("failed to read feature from associated feature file. %v", err)
	}

	if feature.CustomVersionPrefix != config.AppConfig().TagPrefix() && feature.CustomVersionPrefix != "" {
		logging.Instance().Debugf("setting application preset for prefix: %s", feature.CustomVersionPrefix)
		config.AppConfig().SetTagPrefix(feature.CustomVersionPrefix)
	}

	r, err := git.NewRepository()
	if err != nil {
		return err
	}

	*r.FeatureBranch = *r.CurrentBranch

	if r.VersionPrefix != config.AppConfig().TagPrefix() {
		logging.Instance().Warnf("feature version prefix specified does not match existing prefix for this git project ('%s' != '%s')", config.AppConfig().TagPrefix(), r.VersionPrefix)
		if c := prompt.String("continue with feature release (Y/n)? "); strings.ToUpper(c) != "Y" {
			logging.Instance().Info("safely exiting feature release")
			return nil
		}
		logging.Instance().Info("continuing with feature release against warning")
	}

	if err := r.PullChanges(); err != nil {
		return fmt.Errorf("failed to ensure %s is up to date with remote. %v", r.CurrentBranch, err)
	}

	updatedVersion := bumpReleaseVersion(r.LastTag, fc.action)

	if !fc.noChangelog && !fc.noTag {
		changelogEntry := changelog.NewChangelogEntry(feature, r, updatedVersion, fc.action == "MAJOR" || fc.action == "MINOR")
		changelogLines, err := changelog.CreateChangeLogLines(changelogEntry)
		if err != nil {
			return fmt.Errorf("failed to update the changelog. %v", err)
		}

		if err := changelog.WriteChangelogToFile(changelogLines); err != nil {
			return fmt.Errorf("failed to write changelog entry. %v", err)
		}
	}

	if err := os.RemoveAll(GOGDir); err != nil {
		return fmt.Errorf("failed to remove GOG directory. %v", err)
	}

	if err := r.StageChanges(); err != nil {
		return fmt.Errorf("failed to stage removal of GOG metadata folder on %s. %v", r.CurrentBranch, err)
	}

	if err := r.CommitChanges("remove GOG metadata folder"); err != nil {
		return err
	}

	if err := r.Rebase(); err != nil {
		return fmt.Errorf("failed to rebase commits into new release. %v", err)
	}

	if err := r.CheckoutBranch(r.DefaultBranch, false, false); err != nil {
		return fmt.Errorf("failed to checkout branch (%s). %v", r.DefaultBranch, err)
	}

	if err := r.PullChanges(); err != nil {
		return fmt.Errorf("failed to ensure %s is up to date with remote. %v", r.CurrentBranch, err)
	}

	if err := r.SquashMerge(); err != nil {
		return fmt.Errorf("failed to perform squash-merge for new release. %v", err)
	}

	if err := r.StageChanges(); err != nil {
		return fmt.Errorf("failed to stage final changes to %s. %v", r.CurrentBranch, err)
	}

	if err := r.CommitChanges(strings.Join([]string{feature.Jira, feature.Comment}, " ")); err != nil {
		return fmt.Errorf("failed to commit final changes to %s. %v", r.CurrentBranch, err)
	}

	if err := r.Push(); err != nil {
		return fmt.Errorf("failed to push final changes to %s. %v", r.CurrentBranch, err)
	}

	if !fc.noTag {
		if err := feature.CreateReleaseTags(r, updatedVersion); err != nil {
			return fmt.Errorf("failed to create release tags. %v", err)
		}

		if err := r.PushTags(); err != nil {
			return fmt.Errorf("failed to publish release tags to remote. %v", err)
		}
	}

	if err := r.DeleteBranch(r.FeatureBranch); err != nil {
		return fmt.Errorf("failed to delete existing feature branch for %s. %v", feature.Jira, err)
	}

	logging.Instance().Infof("Successfully created new feature release for %s!", feature.Jira)

	return nil
}

func (fc *FinishCommand) Name() string {
	return fc.name
}

func (fc *FinishCommand) Alias() string {
	return fc.alias
}