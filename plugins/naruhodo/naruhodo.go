package naruhodo

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
)

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

type Plugin struct {
}

func (r Plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return strings.Index(message, "なるほど") != -1, message
}

func (r Plugin) DoAction(event plugins.BotEvent, message string) bool {
	idx := int(rd.Int() % len(naruhodoMap))
	event.Reply(naruhodoMap[idx])
	return false // next ng
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)
