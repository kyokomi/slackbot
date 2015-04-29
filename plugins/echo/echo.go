package echo

import (
	"github.com/kyokomi/slackbot/plugins"
	"github.com/nlopes/slack"
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

func (r EchoMessage) DoAction(ctx context.Context, msEvent *slack.MessageEvent, message string, sendMessageFunc func(message string)) bool {
	sendMessageFunc(message)
	return true // next ok
}

var _ plugins.BotMessagePlugin = (*EchoMessage)(nil)
