package main

import (
	"errors"
	"fmt"
	"os"

	"sykesdev.ca/gog/command"
)

func root() error {
	if len(os.Args[1:]) < 1 {
		return errors.New("you must pass a sub-command")
	}

	cmds := []command.Runnable {
		command.NewFeatureCommand(),
		command.NewPushCommand(),
		command.NewFinishCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
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