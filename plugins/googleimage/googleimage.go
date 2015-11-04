package googleimage

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/kyokomi/slackbot/plugins"
)

const endpointURL = "https://www.googleapis.com/customsearch/v1"

type plugin struct {
	rd     *rand.Rand
	client *http.Client
	cx     string
	apiKey string
}

func NewPlugin(cx string, apiKey string) plugins.BotMessagePlugin {
	return &plugin{
		rd:     rand.New(rand.NewSource(time.Now().UnixNano())),
		client: http.DefaultClient,
		cx:     cx,
		apiKey: apiKey,
	}
}

func (p *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "image me")
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	query := strings.Replace(strings.TrimLeft(message, "image me"), "image me", "", 1)

	params := url.Values{}
	params.Set("searchType", "image")
	params.Set("alt", "json")
	params.Set("cx", p.cx)
	params.Set("key", p.apiKey)
	params.Set("q", query)

	resp, err := p.client.Get(fmt.Sprintf("%s?%s", endpointURL, params.Encode()))
	if err != nil {
		log.Println(err)
		return false
	}

	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	items, err := j.Get("items").Array()
	if err != nil {
		log.Println(err)
		return false
	}

	var links []string
	for _, item := range items {
		link := item.(map[string]interface{})["link"].(string)
		links = append(links, link)
	}

	idx := int(p.rd.Int() % len(links))
	event.Reply(links[idx])

	return false // next ng
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
