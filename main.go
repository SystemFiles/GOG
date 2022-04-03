package main

import (
	"errors"
	"fmt"
	"os"

	"sykesdev.ca/gog/command"
	"sykesdev.ca/gog/lib"
)

var Version string

func root() error {
	if len(os.Args[1:]) < 1 {
		return errors.New("you must pass a sub-command\nUsage: gog <feature | push | finish> [options ...] [-h] [-help]")
	}

	cmds := []command.Runnable {
		command.NewFeatureCommand(),
		command.NewPushCommand(),
		command.NewFinishCommand(),
		command.NewUpdateSelfCommand(),
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

	// updater tests
	// u := update.NewUpdater("")
	// a, err := u.GetLatestReleaseAsset()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Println(a)

	// fmt.Println(runtime.GOOS)
	// fmt.Println(runtime.GOARCH)

	// downloadUrl := "https://github.com/SystemFiles/stay-up/releases/download/v2.0.0/stay-up_2.0.0_darwin_amd64.tar.gz"

	// resp, err := http.Get(downloadUrl)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer resp.Body.Close()

	// out, err := os.Create("dist/stay-up_2.0.0_darwin_amd64.tar.gz")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer out.Close()

	// _, err = io.Copy(out, resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	if err := root(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}