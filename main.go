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
			featureCmd := flag.NewFlagSet("feature", flag.ExitOnError)
			fromFeature := featureCmd.Bool("from-feature", false, "-from-feature specifies if this feature will be based on the a current feature branch")
			featureCmd.Parse(os.Args[2:])

			command.ExecFeature(*fromFeature)
		case subcommands[1]:
			command.ExecPush()
		case subcommands[2]:
			finishCmd := flag.NewFlagSet("finish", flag.ExitOnError)

			major := finishCmd.Bool("major", false, "-major specifies that this is a major feature (breaking changes)")
			minor := finishCmd.Bool("minor", false, "-minor specifies that this is a minor feature (no breaking, but is not a bug fix or patch)")
			patch := finishCmd.Bool("patch", false, "-patch specifies this is a bugfix or small patch/update")

			finishCmd.Parse(os.Args[2:])

			command.ExecFinish(*major, *minor, *patch)
		default:
			lib.GetLogger().Error("Invalid subcommand specified")
			help()
			os.Exit(1)
		}
	}
}