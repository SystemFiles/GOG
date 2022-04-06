package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"sykesdev.ca/gog/common"
	"sykesdev.ca/gog/git"
	"sykesdev.ca/gog/logging"
	"sykesdev.ca/gog/models"
)

func FeatureUsage() {
	logging.GetLogger().Info("Usage: gog feature <jira_name> <comment> [-from-feature]")
}

type FeatureCommand struct {
	fs *flag.FlagSet
	
	name string
	alias string
	Jira string
	Comment string
	FromFeature bool
}

func NewFeatureCommand() *FeatureCommand {
	fc := &FeatureCommand{
		name: "feature",
		alias: "feat",
		fs: flag.NewFlagSet("feature", flag.ContinueOnError),
	}

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
	validJiraFormat, err := regexp.Match(`^[A-Z].*\-[0-9].*$`, []byte(fc.Jira))
	if err != nil {
		return fmt.Errorf("failed to parse regular expression for Jira format. %v", err)
	}

	if !validJiraFormat {
		return errors.New("invalid Jira format ... example of a valid format would be 'JIRA-0023'")
	}

	workingDir, GOGDir := common.WorkspacePaths()

	if !git.IsValidRepo() {
		return fmt.Errorf("the current directory (%s) is not a valid git repository", workingDir)
	}

	feature, err := models.NewFeature(fc.Jira, fc.Comment)
	if err != nil {
		return fmt.Errorf("failed to create feature object. %v", err)
	}

	initial_branch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch of git repository. %v", err)
	}

	if feature.BranchExists() {
		return fmt.Errorf("there is already a branch in this repo named %s", feature.Jira)
	}

	if git.HasUnstagedCommits() {
		return fmt.Errorf("there is unstaged commits on your current branch (%s). For your safety, please stage or discard the changes to continue", initial_branch)
	}

	if git.HasUncommittedChanges() {
		return fmt.Errorf("there are staged commits on your current branch (%s) which have not been committed. %v", initial_branch, err)
	}

	if !fc.FromFeature {
		if stderr, err := git.CheckoutDefaultBranch(); err != nil {
			return fmt.Errorf("failed to checkout default branch for repo. %v\n%s", err, stderr)
		}
	}

	if common.PathExists(GOGDir) {
		return fmt.Errorf("%s already exists ... there could already be a feature here. Please fix this and try again", GOGDir)
	}

	if stderr, err := feature.CreateBranch(); err != nil {
		return fmt.Errorf("failed to create or checkout new feature branch, %s. %v \n%s", feature.Jira, err, stderr)
	}

	if err := feature.Save(); err != nil {
		if err := feature.Clean(); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to exit cleanly ... %v", err))
		}

		return fmt.Errorf("failed to create feature tracking file (%v)", err)
	}

	logging.GetLogger().Info(fmt.Sprintf("Successfully created feature %s!", feature.Jira))

	return nil
}

func (fc *FeatureCommand) Name() string {
	return fc.name
}

func (fc *FeatureCommand) Alias() string {
	return fc.alias
}