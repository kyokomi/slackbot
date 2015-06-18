package naruhodo

import (
	"strings"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("naruhodoMessage"), NaruhodoMessage{})
}

type NaruhodoMessage struct {
}

func (r NaruhodoMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	return strings.Index(message, "なるほど") != -1, message
}

func (r NaruhodoMessage) DoAction(ctx context.Context, message string) bool {
	plugins.SendMessage(ctx, "なるほどなるほどですぞ!")
	return false // next ng
}

var _ plugins.BotMessagePlugin = (*NaruhodoMessage)(nil)
