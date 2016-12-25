slackbot [![Coverage Status](https://coveralls.io/repos/kyokomi/slackbot/badge.svg?branch=master&service=github)](https://coveralls.io/github/kyokomi/slackbot?branch=master)
========================

[![wercker status](https://app.wercker.com/status/f609e74d43a232c26e011144646c2cf4/m/master "wercker status")](https://app.wercker.com/project/byKey/f609e74d43a232c26e011144646c2cf4)

Plugin extension a simple slack bot for golang.

## Description

`plugins.BotMessagePlugin`を実装して、`slackbot.BotContext`に`AddPlugin`するだけでプラグイン追加できます。 

Bot側の実装は、[こちら](https://github.com/kyokomi/nepu-bot/blob/master/main.go)を参考にしてください。

## Usage

```go
package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins/echo"
)

func main() {
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.Parse()

	bot, err := slackbot.NewBotContext(token)
	if err != nil {
		panic(err)
	}
	bot.AddPlugin("echo", echo.NewPlugin())

	bot.WebSocketRTM()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":8000", nil)
}
```

## Plugin Example

```go
package echo

import (
	"github.com/kyokomi/slackbot/plugins"
)

type plugin struct {
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{}
}

func (p *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return true, message
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	event.Reply(message)
	return true // next ok
}

func (p *plugin) Help() string {
	return `echo:
	all message echo
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
``` 

## License

[MIT](https://github.com/kyokomi/slackbot/blob/master/LICENSE)
