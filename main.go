package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

		// create options
		jira := flag.Arg(1)
		comment := strings.Join(flag.Args()[2:], " ")
		fromFeature := *flag.Bool("from-feature", false, "specifies if this feature will be based on the a current feature branch")

		switch subCmd {
		case subcommands[0]:
			command.ExecFeature(jira, comment, fromFeature)
		case subcommands[1]:
			return
		case subcommands[2]:
			return
		default:
			lib.GetLogger().Error("Invalid subcommand specified")
			help()
			os.Exit(1)
		}
	}
}