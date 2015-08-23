slackbot
========================

[![Circle CI](https://circleci.com/gh/kyokomi/slackbot.svg?style=svg)](https://circleci.com/gh/kyokomi/slackbot)
[![Coverage Status](https://coveralls.io/repos/kyokomi/slackbot/badge.svg?branch=master&service=github)](https://coveralls.io/github/kyokomi/slackbot?branch=master)

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
	bot.AddPlugin("echo", echo.EchoMessage{})

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

type EchoMessage struct {
}

func (r EchoMessage) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return true, message
}

func (r EchoMessage) DoAction(event plugins.BotEvent, message string) bool {
	event.Reply(message)
	return true // next ok
}

var _ plugins.BotMessagePlugin = (*EchoMessage)(nil)
``` 

## License

[MIT](https://github.com/kyokomi/slackbot/blob/master/LICENSE)
