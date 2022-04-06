package main

import (
	"errors"
	"fmt"
	"os"

	"sykesdev.ca/gog/command"
	"sykesdev.ca/gog/common"
	"sykesdev.ca/gog/config"
	"sykesdev.ca/gog/logging"
	"sykesdev.ca/gog/update"
)

func root() error {
	if len(os.Args[1:]) < 1 {
		return errors.New("you must pass a sub-command\nUsage: gog <feature | push | finish> [options ...] [-h] [-help]")
	}

	if common.StringInSlice(os.Args, "-v") || common.StringInSlice(os.Args, "-version") {
		logging.GetLogger().Info(fmt.Sprintf("Current Version of GOG: v%s", update.Version))
		return nil
	}

	cmds := []command.Runnable {
		command.NewFeatureCommand(),
		command.NewPushCommand(),
		command.NewFinishCommand(),
		command.NewUpdateSelfCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand || cmd.Alias() == subcommand {
			if common.StringInSlice(os.Args, "-h") || common.StringInSlice(os.Args, "-help") {
				cmd.Help()
				return nil
			}
			if err := cmd.Init(os.Args[2:]); err != nil {
				return err
			}
			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}

func main() {
	logging.GetLogger().SetupLogger(config.AppConfig().LogLevel())

	if err := root(); err != nil {
		logging.GetLogger().Error(err.Error())
		os.Exit(1)
	}
}