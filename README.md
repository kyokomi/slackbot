slackbot
========================

[![Circle CI](https://circleci.com/gh/kyokomi/slackbot.svg?style=svg)](https://circleci.com/gh/kyokomi/slackbot)

Plugin extension a simple slack bot for golang.

## Description

`plugins.BotMessagePlugin`を実装して、`init()`でplugings.AddPluginを呼び出すだけでbotのプラグインを追加できます。

Bot側の実装は、[こちら](https://github.com/kyokomi/nepu-bot/blob/master/main.go)を参考にしてください。

## Example

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

