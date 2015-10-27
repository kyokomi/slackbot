package lgtm

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
}

const lgtmURL = "http://lgtm.in/g"

func (m Plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "lgtm")
}

func (m Plugin) DoAction(event plugins.BotEvent, message string) bool {
	randomURL := lgtmURL
	args := strings.Fields(message)
	if len(args) == 2 {
		randomURL += "/" + args[1]
	}
	sendMessage, isNext := GetLGTMImageURL(randomURL)

	event.Reply(sendMessage)

	return isNext // next stop
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)

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
