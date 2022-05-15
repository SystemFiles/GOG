package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"sykesdev.ca/gog/common"
	"sykesdev.ca/gog/git"
	"sykesdev.ca/gog/logging"
	"sykesdev.ca/gog/models"
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
	if !git.IsValidRepo() {
		return fmt.Errorf("the current directory is not a valid git repository")
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
	
	if pc.message == "" {
		pc.message = fmt.Sprintf("Test Build (%d)", feature.TestCount)
		feature.UpdateTestCount()
	}

	if stderr, err := feature.PushChanges(pc.message); err != nil {
		return fmt.Errorf("failed to push changes to remote repository. %v \n%s", err, stderr)
	}

	logging.GetLogger().Info("Successfully pushed changes to remote feature!")

	return nil
}

func (pc *PushCommand) Name() string {
	return pc.name
}

func (pc *PushCommand) Alias() string {
	return pc.alias
}