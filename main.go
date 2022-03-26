package main

import (
	"errors"
	"fmt"
	"os"

	"sykesdev.ca/gog/command"
	"sykesdev.ca/gog/lib"
)

func root() error {
	if len(os.Args[1:]) < 1 {
		return errors.New("you must pass a sub-command\nUsage: gog <feature | push | finish> [options ...] [-h] [-help]")
	}

	cmds := []command.Runnable {
		command.NewFeatureCommand(),
		command.NewPushCommand(),
		command.NewFinishCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			if lib.StringInSlice(os.Args, "-h") || lib.StringInSlice(os.Args, "-help") {
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
	if err := root(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}