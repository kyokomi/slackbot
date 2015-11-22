package googleimage

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
)

const endpointURL = "https://www.googleapis.com/customsearch/v1"

type plugin struct {
	rd     *rand.Rand
	client GoogleImageAPIClient
	cx     string
	apiKey string
}

func NewPlugin(client GoogleImageAPIClient) plugins.BotMessagePlugin {
	return &plugin{
		rd:     rand.New(rand.NewSource(time.Now().UnixNano())),
		client: client,
	}
}

func (p *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "image me")
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	query := strings.Replace(strings.TrimLeft(message, "image me"), "image me", "", 1)

	links, err := p.client.GetImageLinks(query)
	if err != nil {
		log.Println(err)
		return false
	}

	idx := int(p.rd.Int() % len(links))
	event.Reply(links[idx])

	return false // next ng
}

func (p *plugin) Help() string {
	return `googleimage: グーグル画像検索
	image me [検索キーワード]
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
