package main

import (
	"errors"
	"fmt"
	"os"

	"sykesdev.ca/gog/cmd"
	"sykesdev.ca/gog/config"
	"sykesdev.ca/gog/internal/common"
	"sykesdev.ca/gog/internal/logging"
	"sykesdev.ca/gog/internal/update"
)

func root() error {
	if len(os.Args[1:]) < 1 {
		return errors.New("you must pass a sub-command\nUsage: gog <feature | push | finish> [options ...] [-h] [-help]")
	}

	if common.StringInSlice(os.Args, "-v") || common.StringInSlice(os.Args, "-version") {
		logging.GetLogger().Info(fmt.Sprintf("Current Version of GOG: %s", update.Version))
		return nil
	}

	cmds := []cmd.Runnable {
		cmd.NewFeatureCommand(),
		cmd.NewPushCommand(),
		cmd.NewFinishCommand(),
		cmd.NewUpdateSelfCommand(),
		cmd.NewSimplePushCommand(),
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