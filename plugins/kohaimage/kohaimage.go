package kohaimage

import (
	"github.com/kyokomi/slackbot/plugins"
)

type plugin struct {
	kohaAPI KohaAPI
}

func NewPlugin(kohaAPI KohaAPI) plugins.BotMessagePlugin {
	return &plugin{kohaAPI: kohaAPI}
}

func (r plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "koha")
}

func (r plugin) DoAction(event plugins.BotEvent, message string) bool {
	event.Reply(r.kohaAPI.GetImageURL())
	return false // next ng
}

func (p *plugin) Help() string {
	return `kohaimage: コハエース画像
	文中にkohaが含まれていると、経験値先生力作のLINEスタンプ申請落ちした画像をランダムで表示します。
`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
