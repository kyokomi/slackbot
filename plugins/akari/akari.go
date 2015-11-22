package akari

import (
	"fmt"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
)

const keyword = "大好き"

type plugin struct {
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{}
}

func (p *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeywords(message, keyword)
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	daisuki := strings.Replace(strings.TrimLeft(message, keyword), keyword, "", 1)
	event.Reply(fmt.Sprintf("わぁい%s あかり%s大好き", daisuki, daisuki))
	return false // next ok
}

func (p *plugin) Help() string {
	return `あかり:
	大好きに反応する わぁいXXXX あかりXXXX大好き
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
