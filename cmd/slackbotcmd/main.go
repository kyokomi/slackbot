package main

import (
	"os"

	"github.com/urfave/cli"
)

//go:generate ego -o=command/ego.go -package=command -version=false command/templates

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "kyokomi"
	app.Email = ""
	app.Usage = ""

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}
