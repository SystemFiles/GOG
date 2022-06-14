package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/git"
	"sykesdev.ca/gog/internal/logging"
	"sykesdev.ca/gog/internal/models"
)

type PushCommand struct {
	fs *flag.FlagSet

	name string
	alias string
	message string
}

func NewPushCommand() *PushCommand {
	pc := &PushCommand{
		name: "push",
		alias: "p",
		fs: flag.NewFlagSet("push", flag.ContinueOnError),
	}

	pc.fs.Usage = pc.Help

	return pc
}

func (fc *PushCommand) Help() {
	fmt.Printf(
`Usage: %s (%s | %s) [message] [-h] [-help]

-------====== Push Arguments ======-------

message
	specifies a commit message for this feature push
`, os.Args[0], fc.name, fc.alias)

	fc.fs.PrintDefaults()

	fmt.Println("\n-------================================-------")
}

func (pc *PushCommand) Init(args []string) error {
	err := pc.fs.Parse(args)

	if len(pc.fs.Args()) >= 1 {
		pc.message = strings.Join(pc.fs.Args(), " ")
	}

	return err
}

func (pc *PushCommand) Run() error {
	r, err := git.NewRepository()
	if err != nil {
		return err
	}

	GOGDir := common.GOGPath()

	if !common.PathExists(GOGDir + "/feature.json") {
		return errors.New("feature file not found ... there may not be a GOG feature on this branch")
	}

	feature, err := models.NewFeatureFromFile()
	if err != nil {
		return fmt.Errorf("failed to read feature from features file (%s). %v", GOGDir + "/feature.json", err)
	}
	defer feature.Save()
	
	r.FeatureBranch = r.CurrentBranch

	if pc.message == "" {
		pc.message = fmt.Sprintf("%s Test Build (%d)", feature.Jira, feature.TestCount)
		feature.UpdateTestCount()
	} else {
		pc.message = strings.Join([]string{feature.Jira, pc.message}, " ")
	}

	if err := r.StageChanges(); err != nil {
		return fmt.Errorf("failed to stage current changes for %s. %v", r.CurrentBranch, err)
	}

	if err := r.CommitChanges(pc.message); err != nil {
		return fmt.Errorf("failed to commit current changes for %s. %v", r.CurrentBranch, err)
	}

	if r.CurrentBranch.RemoteExists {
		if err := r.PullChanges(); err != nil {
			return fmt.Errorf("failed to ensure %s is up to date with remote. %v", r.CurrentBranch, err)
		}
	}

	if err := r.Push(); err != nil {
		return fmt.Errorf("failed to push local commits to remote. %v", err)
	}

	logging.Instance().Info("Successfully pushed changes to remote feature!")

	return nil
}

func (pc *PushCommand) Name() string {
	return pc.name
}

func (pc *PushCommand) Alias() string {
	return pc.alias
}