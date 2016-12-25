package main

import (
	"fmt"
	"os"

	"github.com/kyokomi/slackbot/cmd/slackbotcmd/command"
	"github.com/urfave/cli"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "new",
		Usage:  "",
		Action: command.CmdNew,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "pkg", Value: "", Usage: "package name", EnvVar: ""},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
