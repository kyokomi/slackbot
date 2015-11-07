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
