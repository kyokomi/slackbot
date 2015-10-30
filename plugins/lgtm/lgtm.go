package lgtm

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kyokomi/slackbot/plugins"
)

const lgtmURL = "http://lgtm.in/g"

type plugin struct {
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{}
}

func (p *plugin) buildRandomURL(message string) string {
	randomURL := lgtmURL
	args := strings.Fields(message)
	if len(args) == 2 {
		randomURL += "/" + args[1]
	}
	return randomURL
}

func (p *plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "lgtm")
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	sendMessage, isNext := GetLGTMImageURL(p.buildRandomURL(message))

	event.Reply(sendMessage)

	return isNext // next stop
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)

func GetLGTMImageURL(lgtmURL string) (string, bool) {
	res, err := http.Get(lgtmURL)
	if err != nil {
		return err.Error(), true
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return err.Error(), true
	}

	text, exists := doc.Find("#imageUrl").Attr("value")
	if !exists {
		return lgtmURL + ": not exists", true
	}

	return text, false
}
