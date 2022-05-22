package cmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/git"
	"sykesdev.ca/gog/internal/logging"
	"sykesdev.ca/gog/internal/prompt"
)

type SimplePushCommand struct {
	fs *flag.FlagSet

	name string
	alias string
	message string
}

func gitLatestTestBuild() int {
	prev, err := git.GetPreviousNCommitMessage(1)
	if err != nil {
		return -1
	}

	re, err := regexp.Compile(`\([0-9]\)$`)
	if err != nil {
		return -1
	}

	if loc := re.FindStringIndex(prev[0]); loc != nil {
		res, err := strconv.ParseInt(string(prev[0][loc[0]+1]), 10, 32)
		if err != nil {
			return -1
		}

		return int(res)
	}

	return -1
}

func NewSimplePushCommand() *SimplePushCommand {
	c := &SimplePushCommand{
		name: "simple-push",
		alias: "sp",
		fs: flag.NewFlagSet("simple-push", flag.ContinueOnError),
	}

	c.fs.Usage = c.Help

	return c
}

func (c *SimplePushCommand) Init(args []string) error {
	err := c.fs.Parse(args)
	if err != nil {
		return err
	}

	// optional message if user wants custom
	if len(c.fs.Args()) >= 1 {
		c.message = strings.Join(c.fs.Args(), " ")
	}

	return nil
}

func (c *SimplePushCommand) Help() {
	fmt.Printf(
`Usage: %s %s [message] [-h] [-help]

Simple-Push is a utility to allow non-feature related code pushes directly to the current remote branch. If used without a message one will be generated.

-------====== Simple-Push Arguments ======-------

message
	(optionally) specifies a commit message for this simple push operation
`, os.Args[0], c.name)
	
	c.fs.PrintDefaults()

	fmt.Println("\n-------================================-------")
}

func (c *SimplePushCommand) Run() error {
	if !git.IsValidRepo() {
		return fmt.Errorf("the current directory is not a valid git repository")
	}

	GOGDir := common.GOGPath()

	if common.PathExists(GOGDir + "/feature.json") {
		logging.Instance().Warn("this project seems to already have an associated GOG feature file. It is recommended to use 'gog (p | push)' for feature code changes")
		contResp := prompt.String("would you like to continue? [Y/N] ")

		if strings.ToUpper(contResp) != "Y" {
			logging.Instance().Info("cancelling simple-push ...")
			return nil
		}
	}

	cbOut, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failure to get name for current branch. %v\n%s", err, cbOut)
	}

	buildNumber := gitLatestTestBuild()

	if c.message == "" {
		c.message = fmt.Sprintf("Test Build (%d)", buildNumber + 1)
	}	else {
		c.message = fmt.Sprintf("%s (%d)", c.message, buildNumber + 1)
	}

	if stderr, err := git.StageChanges(); err != nil {
		return fmt.Errorf("failure to stage local changes. %v\n%s", err, stderr)
	}

	var pushArgs string
	if git.HasUncommittedChanges() {
		formattedMessage := fmt.Sprintf("%s %s", cbOut, c.message)
		if stderr, err := git.Commit(formattedMessage); err != nil {
			return fmt.Errorf("failed to commit local changes. %v\n%s", err, stderr)
		}
	}

	if !git.RemoteBranchExists(cbOut) {
		pushArgs = fmt.Sprintf("--set-upstream origin %s", cbOut)
	} else {
		// only pull changes if a remote exists
		if stderr, err := git.PullChanges(); err != nil {
			return fmt.Errorf("failed to pull changes from existing remote branch. %v\n%s", err, stderr)
		}
	}

	if stderr, err := git.PushRemote(pushArgs); err != nil {
		return fmt.Errorf("failed to push local changes to remote. %v\n%s", err, stderr)
	}

	logging.Instance().Info("Successfully pushed changes to remote (" + cbOut + ")!")
	
	return nil
}

func (c *SimplePushCommand) Name() string {
	return c.name
}

func (c *SimplePushCommand) Alias() string {
	return c.alias
}
