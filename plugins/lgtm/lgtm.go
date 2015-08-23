package lgtm

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
}

const lgtmURL = "http://www.lgtm.in/g"

func (m Plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "lgtm")
}

func (m Plugin) DoAction(event plugins.BotEvent, message string) bool {
	sendMessage, isNext := GetLGTMImageURL(lgtmURL)

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
