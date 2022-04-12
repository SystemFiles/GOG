package command

import (
	"flag"
	"fmt"
	"os"

	"sykesdev.ca/gog/logging"
	"sykesdev.ca/gog/update"
)

type UpdateSelfCommand struct {
	fs *flag.FlagSet

	name string
	tag string
}

func NewUpdateSelfCommand() *UpdateSelfCommand {
	usc := &UpdateSelfCommand{
		name: "update",
		fs: flag.NewFlagSet("update", flag.ContinueOnError),
	}

	usc.fs.StringVar(&usc.tag, "tag", "", "specifies a specific version tag to use for update")

	return usc
}

func (usc *UpdateSelfCommand) Help() {
	fmt.Printf(
`Usage: %s update [-tag TAG] [-h] [-help]

--------======= Tag Arguments =======--------

`, os.Args[0])

	usc.fs.PrintDefaults()

	fmt.Println("\n-------================================-------")
}

func (usc *UpdateSelfCommand) Init(args []string) error {
	return usc.fs.Parse(args)
}

func (usc *UpdateSelfCommand) Run() error {
	logging.GetLogger().Info("Performing in-place upgrade for GOG ...")

	u, err := update.NewUpdater(usc.tag)
	if err != nil {
		return err
	}

	if err := u.Update(); err != nil {
		return err
	}

	logging.GetLogger().Info(fmt.Sprintf("Successfully updated GOG from %s to %s", u.CurrentVersion(), u.UpdateVersion()))
	return nil
}

func (usc *UpdateSelfCommand) Name() string {
	return usc.name
}

func (usc *UpdateSelfCommand) Alias() string {
	return ""
}