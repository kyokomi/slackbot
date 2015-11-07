package suddendeath

import (
	"strings"
	"unicode/utf8"

	"github.com/kyokomi/slackbot/plugins"
)

type plugin struct {
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{}
}

func (r *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return strings.Index(message, "突然の") != -1, message
}

func (r *plugin) DoAction(event plugins.BotEvent, message string) bool {
	size := utf8.RuneCountInString(message)
	header := ""
	for i := 0; i < size+2; i++ {
		header += "人"
	}

	fotter := ""
	for i := 0; i < size; i++ {
		fotter += "^Y"
	}

	reMessage := "＿" + header + "＿"
	reMessage += "\n"
	reMessage += "＞　" + message + "　＜"
	reMessage += "\n"
	reMessage += "￣Y" + fotter + "￣"

	event.Reply(reMessage)

	return false // next ok
}

func (p *plugin) Help() string {
	return `suddendeath: 突然の死

	突然の<free message>:

		＿人人人人人人人人人人人人＿
		＞　突然のfree message　＜
		￣Y^Y^Y^Y^Y^Y^Y^Y^Y^Y￣
`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
