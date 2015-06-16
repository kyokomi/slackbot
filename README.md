slackbot
========================

[![Circle CI](https://circleci.com/gh/kyokomi/slackbot.svg?style=svg)](https://circleci.com/gh/kyokomi/slackbot)

Plugin extension a simple slack bot for golang.

## Description

`plugins.BotMessagePlugin`を実装して、`init()`でplugings.AddPluginを呼び出すだけでbotのプラグインを追加できます。

Bot側の実装は、[こちら](https://github.com/kyokomi/nepu-bot/blob/master/main.go)を参考にしてください。

## Usage

```go
package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins"

	_ "github.com/kyokomi/slackbot/plugins/echo"
)

func main() {
	var token string
	flag.StringVar(&token, "token", os.Getenv("SLACK_BOT_TOKEN"), "SlackのBotToken")
	flag.Parse()

	ctx := plugins.Context()

	c := slackbot.DefaultConfig()
	c.Name = "bot name"
	c.SlackToken = token

	slackbot.WebSocketRTM(ctx, c)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})
	http.ListenAndServe(":8000", nil)
}
```

## Plugin Example

```go
package echo

import (
	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("echoMessage"), EchoMessage{})
}

type EchoMessage struct {
}

func (r EchoMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	return true, message
}

func (r EchoMessage) DoAction(ctx context.Context, message string) bool {
	plugins.SendMessage(ctx, message)
	return true // next ok
}

var _ plugins.BotMessagePlugin = (*EchoMessage)(nil)
```

