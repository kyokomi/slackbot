package suddendeath

import (
	"strings"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("suddendeath"), SuddenDeathMessage{})
}

type SuddenDeathMessage struct {
}

func (r SuddenDeathMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	return strings.Index(message, "突然の") != -1, message
}

func (r SuddenDeathMessage) DoAction(ctx context.Context, message string) bool {
	header := ""
	for i := 0; i < len(message)/2; i++ {
		header += "人"
	}

	fotter := ""
	for i := 0; i < len(message)/2; i++ {
		fotter += "YY"
	}

	reMessage := "＿" + header + "＿"
	reMessage += "\n"
	reMessage += "＞　" + message + "　＜"
	reMessage += "\n"
	reMessage += "￣" + fotter + "￣"

	plugins.SendMessage(ctx, reMessage)
	return false // next ok
}

var _ plugins.BotMessagePlugin = (*SuddenDeathMessage)(nil)
