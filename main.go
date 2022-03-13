package main

import (
	"flag"
	"fmt"
	"os"

	"sykesdev.ca/gog/command"
	"sykesdev.ca/gog/lib"
)

var subcommands = [3]string{"feature", "push", "finish"}

func help() {
	lib.GetLogger().Info("Usage: gog <subCmd> [options...]")
	lib.GetLogger().Info("Sub Commands Available:")
	for _, s := range subcommands {
		lib.GetLogger().Info(fmt.Sprintf(" -- %s", s))
	}
}

func main() {
	flag.Parse()
	
	if len(flag.Args()) == 0 {
		lib.GetLogger().Error("Incorrect number of arguments passed")
		help()
		os.Exit(0)
	}

	if len(flag.Args()) >= 1 {
		subCmd := flag.Arg(0)

		switch subCmd {
		case subcommands[0]:
			command.ExecFeature()
		case subcommands[1]:
			command.ExecPush()
		case subcommands[2]:
			return
		default:
			lib.GetLogger().Error("Invalid subcommand specified")
			help()
			os.Exit(1)
		}
	}
}