package naruhodo

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
)

type plugin struct {
	rd          *rand.Rand
	naruhodoMap []string
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{
		rd: rand.New(rand.NewSource(time.Now().UnixNano())),
		naruhodoMap: []string{
			"なるほどなるほどですぞ!",
			"なるほど!",
			"なるほど?",
			"なーるほど！",
			"それはなるほどですね",
			"なるほど!!",
			"なるほど!!!",
		},
	}
}

func (p *plugin) getRandomNaruhodo() string {
	idx := int(p.rd.Int() % len(p.naruhodoMap))
	return p.naruhodoMap[idx]
}

func (p *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return strings.Index(message, "なるほど") != -1, message
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	event.Reply(p.getRandomNaruhodo())
	return false // next ng
}

func (p *plugin) Help() string {
	return `naruhodo: なるほど
	文中に[なるほど]が含まれている、ランダムなるほどメッセージを返します。
`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
