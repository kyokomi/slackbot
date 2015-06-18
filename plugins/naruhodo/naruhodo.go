package naruhodo

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

var naruhodoMap = []string{
	"なるほどなるほどですぞ!",
	"なるほど!",
	"なるほど?",
	"なーるほど！",
	"それはなるほどですね",
	"なるほど!!",
	"なるほど!!!",
}

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	plugins.AddPlugin(pluginKey("naruhodoMessage"), NaruhodoMessage{})
}

type NaruhodoMessage struct {
}

func (r NaruhodoMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	return strings.Index(message, "なるほど") != -1, message
}

func (r NaruhodoMessage) DoAction(ctx context.Context, message string) bool {
	idx := int(rd.Int() % len(naruhodoMap))
	plugins.SendMessage(ctx, naruhodoMap[idx])
	return false // next ng
}

var _ plugins.BotMessagePlugin = (*NaruhodoMessage)(nil)
