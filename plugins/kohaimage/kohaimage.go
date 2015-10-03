package kohaimage

import (
	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
	kohaAPI KohaAPI
}

func NewPlugin(kohaAPI KohaAPI) plugins.BotMessagePlugin {
	return &Plugin{kohaAPI: kohaAPI}
}

func (r Plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "koha")
}

func (r Plugin) DoAction(event plugins.BotEvent, message string) bool {
	event.Reply(r.kohaAPI.GetImageURL())
	return false // next ng
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)
