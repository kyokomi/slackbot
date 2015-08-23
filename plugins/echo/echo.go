package echo

import (
	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
}

func (r Plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return true, message
}

func (r Plugin) DoAction(event plugins.BotEvent, message string) bool {
	event.Reply(message)
	return true // next ok
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)
