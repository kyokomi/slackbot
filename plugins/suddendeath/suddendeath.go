package suddendeath

import (
	"strings"
	"unicode/utf8"

	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
}

func (r Plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return strings.Index(message, "突然の") != -1, message
}

func (r Plugin) DoAction(event plugins.BotEvent, message string) bool {
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

var _ plugins.BotMessagePlugin = (*Plugin)(nil)
