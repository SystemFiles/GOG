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

func gitLatestTestBuild(r *git.Repository) int {
	prev, err := r.LogN(1)
	if err != nil {
		return -1
	}

	re, err := regexp.Compile(`\([0-9]{0,5}\)$`)
	if err != nil {
		return -1
	}

	if match := re.FindStringSubmatch(prev[0]); match != nil {
		res, err := strconv.ParseInt(match[0][1:len(match[0])-1], 10, 64)
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
`Usage: %s (%s | %s) [message] [-h] [-help]

Simple-Push is a utility to allow non-feature related code pushes directly to the current remote branch. If used without a message one will be generated.

-------====== Simple-Push Arguments ======-------

message
	(optionally) specifies a commit message for this simple push operation
`, os.Args[0], c.name, c.alias)
	
	c.fs.PrintDefaults()

	fmt.Println("\n-------================================-------")
}

func (c *SimplePushCommand) Run() error {
	r, err := git.NewRepository()
	if err != nil {
		return err
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

	if c.message == "" {
		buildNumber := gitLatestTestBuild(r)
		c.message = fmt.Sprintf("%s Test Build (%d)", r.CurrentBranch, buildNumber + 1)
	}	else {
		c.message = fmt.Sprintf("%s %s", r.CurrentBranch, c.message)
	}

	if err := r.StageChanges(); err != nil {
		return err
	}

	if err := r.CommitChanges(c.message); err != nil {
		return err
	}
	
	if err := r.PullChanges(); err != nil {
		return err
	}

	if err := r.Push(); err != nil {
		return err
	}

	logging.Instance().Info("Successfully pushed changes to remote (" + r.CurrentBranch.Name + ")!")
	
	return nil
}

func (c *SimplePushCommand) Name() string {
	return c.name
}

func (c *SimplePushCommand) Alias() string {
	return c.alias
}
