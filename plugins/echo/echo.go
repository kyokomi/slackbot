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
