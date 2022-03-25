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
	lib.GetLogger().Info("Usage: gog <sub_command> [options...]")
	lib.GetLogger().Info("Sub Commands Available:")
	for _, s := range subcommands {
		lib.GetLogger().Info(fmt.Sprintf(" -- %s", s))
	}
}

func gogFeature(fromFeature bool) {
	if len(flag.Args()) < 3 {
		lib.GetLogger().Error("Invalid usage of feature sub-command")
		command.FeatureUsage()
		os.Exit(2)
	}

	jira := flag.Arg(1)
	comment := strings.Join(flag.Args()[2:], " ")

	command.ExecFeature(jira, comment, fromFeature)
}

func gogPush() {
	var message string
	if len(flag.Args()) >= 2 {
		message = strings.Join(flag.Args()[1:], " ")
	}

	command.ExecPush(message)
}

func gogFinish(isMajor, isMinor, isPatch bool) {
	if isMajor {
		command.ExecFinish("MAJOR")
	} else if isMinor {
		command.ExecFinish("MINOR")
	} else if isPatch {
		command.ExecFinish("PATCH")
	} else {
		lib.GetLogger().Error("Invalid usage of finish sub-command")
		command.FinishUsage()
		os.Exit(2)
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
			fromFeature := featureCmd.Bool("from-feature", false, "specifies if this feature will be based on the a current feature branch")
			featureCmd.Parse(os.Args[2:])

			gogFeature(*fromFeature)
		case subcommands[1]:
			gogPush()
		case subcommands[2]:
			finishCmd := flag.NewFlagSet("finish", flag.ExitOnError)

			major := finishCmd.Bool("major", false, "specifies that this is a major feature (breaking changes)")
			minor := finishCmd.Bool("minor", false, "specifies that this is a minor feature (no breaking, but is not a bug fix or patch)")
			patch := finishCmd.Bool("patch", false, "specifies this is a bugfix or small patch/update")

			finishCmd.Parse(os.Args[2:])

			gogFinish(*major, *minor, *patch)
		default:
			lib.GetLogger().Error("Invalid subcommand specified")
			help()
			os.Exit(1)
		}
	}
}