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
