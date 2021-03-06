package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"sykesdev.ca/gog/config"
	"sykesdev.ca/gog/internal/git"
	"sykesdev.ca/gog/internal/logging"
	"sykesdev.ca/gog/internal/models"
	"sykesdev.ca/gog/internal/prompt"
)

func FeatureUsage() {
	logging.Instance().Info("Usage: gog feature <jira_name> <comment> [-from-feature]")
}

type FeatureCommand struct {
	fs *flag.FlagSet
	
	name string
	alias string
	Jira string
	Comment string
	CustomVersionPrefix string
	FromFeature bool
}

func NewFeatureCommand() *FeatureCommand {
	fc := &FeatureCommand{
		name: "feature",
		alias: "feat",
		fs: flag.NewFlagSet("feature", flag.ContinueOnError),
	}

	fc.fs.StringVar(&fc.CustomVersionPrefix, "prefix", "", "optionally specifies a version prefix to use for this feature which will override existing prefix in global GOG config")
	fc.fs.BoolVar(&fc.FromFeature, "from-feature", false, "specifies if this feature will be based on the a current feature branch")

	fc.fs.Usage = fc.Help

	return fc
}

func (fc *FeatureCommand) Help() {
	fmt.Printf(
`Usage: %s (%s | %s) <jira> <comment> [-from-feature] [-h] [-help]

-------====== Feature Arguments ======-------

jira
	specifies the JIRA issue we are working under
comment
	specifies a human-readable comment describing the issue/feature

------================================------

`, os.Args[0], fc.name, fc.alias)

	fc.fs.PrintDefaults()

	fmt.Println("\n-------================================-------")
}

func (fc *FeatureCommand) Init(args []string) error {
	err := fc.fs.Parse(args)

	if len(fc.fs.Args()) < 2 {
		return errors.New("invalid usage of feature command. must pass a jira identifier and comment (re-run with -h for full usage details)")
	}

	fc.Jira = fc.fs.Arg(0)
	fc.Comment = strings.Join(fc.fs.Args()[1:], " ")

	return err
}

func (fc *FeatureCommand) Run() error {
	validJiraFormat, err := regexp.Match(`^([A-Za-z])+-([0-9])+$`, []byte(fc.Jira))
	if err != nil {
		return fmt.Errorf("failed to parse regular expression for Jira format. %v", err)
	}

	if !validJiraFormat {
		return errors.New("invalid Jira format ... example of a valid format would be 'JIRA-0023'")
	}

	if fc.CustomVersionPrefix != config.AppConfig().TagPrefix() && fc.CustomVersionPrefix != "" {
		logging.Instance().Debugf("setting application preset for prefix: %s", fc.CustomVersionPrefix)
		config.AppConfig().SetTagPrefix(fc.CustomVersionPrefix)
	}

	r, err := git.NewRepository()
	if err != nil {
		return err
	}

	feature, err := models.NewFeature(fc.Jira, fc.Comment, fc.CustomVersionPrefix)
	if err != nil {
		return fmt.Errorf("failed to create feature object. %v", err)
	}

	if r.VersionPrefix != config.AppConfig().TagPrefix() {
		logging.Instance().Warnf("feature version prefix specified does not match existing prefix for this git project ('%s' != '%s')", config.AppConfig().TagPrefix(), r.VersionPrefix)
		if c := prompt.String("continue with feature creation (Y/n)? "); strings.ToUpper(c) != "Y" {
			logging.Instance().Info("safely exiting feature creation")
			logging.Instance().Info("if you wish to use the existing version prefix, but it is not set in the global config for GOG, you can pass it using the -prefix flag (see -help for details)")
			return nil
		}
		logging.Instance().Info("continuing with feature creation against warning")
	}

	if r.ContainsBranch(feature.Jira) {
		return fmt.Errorf("there is already a branch in this repo named %s", feature.Jira)
	}

	if fc.FromFeature {
		return errors.New("from-feature not implemented yet")
	}

	if err := r.StageChanges(); err != nil {
		return fmt.Errorf("failed to stage changes on current branch (%s) before starting new feature. %v", r.CurrentBranch, err)
	}

	if r.CurrentBranch.UncommittedChanges() && r.CurrentBranch.Name != r.DefaultBranch.Name {
		return fmt.Errorf("the current branch (%s) has uncommitted changes, please review and discard/commit them before starting a new feature", r.CurrentBranch)
	}

	if err := r.CheckoutBranch(r.DefaultBranch, false, false); err != nil {
		return fmt.Errorf("failed to checkout branch %s. %v", r.DefaultBranch, err)
	}

	if err := r.PullChanges(); err != nil {
		return fmt.Errorf("failed to pull some changes before creating the new feature. %v", err)
	}

	r.FeatureBranch = git.NewBranch(feature.Jira)

	if err := r.CheckoutBranch(r.FeatureBranch, true, true); err != nil {
		return fmt.Errorf("failed to create or checkout new feature branch, %s. %v", feature.Jira, err)
	}

	if err := feature.Save(); err != nil {
		return fmt.Errorf("failed to create feature tracking file (%v)", err)
	}

	logging.Instance().Infof("Successfully created feature %s!", feature.Jira)

	return nil
}

func (fc *FeatureCommand) Name() string {
	return fc.name
}

func (fc *FeatureCommand) Alias() string {
	return fc.alias
}