package cmd

import (
	"bee-agent/build"
	"github.com/urfave/cli/v2"
	"os"
)

var rootCmd = &cli.App{
	Name:    "bee-agent",
	Usage:   "A agent for beekeeper.one",
	Version: build.UserVersion(),
	Commands: []*cli.Command{
		runCmd,
	},
}

func Execute() error {
	if err := rootCmd.Run(os.Args); err != nil {
		return err
	}
	return nil
}
