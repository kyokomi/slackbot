# slackbotcmd



## Description

slackbotのplugin開発支援ツール

## Usage

```
$ slackbotcmd --help
NAME:
   slackbotcmd -

USAGE:
   slackbotcmd [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR(S):
   kyokomi

COMMANDS:
   new
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

## Example

```
$ slackbotcmd new --pkg cron
```

Output:

```go
package cron

import (
	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("cronMessage"), CronMessage{})
}

type CronMessage struct {
}

func (r CronMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	// TODO: execute message check
	return true, message
}

func (r CronMessage) DoAction(ctx context.Context, message string) bool {
	// TODO: message action
	//plugins.SendMessage(ctx, message)
	return true // next ok
}

var _ plugins.BotMessagePlugin = (*CronMessage)(nil)
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/kyokomi/slackbot/cmd/slackbotcmd
```

## Contribution

1. Fork ([https://github.com/kyokomi/slackbotcmd/fork](https://github.com/kyokomi/slackbotcmd/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[kyokomi](https://github.com/kyokomi)
