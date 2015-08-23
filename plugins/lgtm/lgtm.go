package lgtm

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
}

func (m Plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "lgtm")
}

func (m Plugin) DoAction(event plugins.BotEvent, message string) bool {
	sendMessage, isNext := getLGTMImageURL()

	event.Reply(sendMessage)

	return isNext // next stop
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)

func getLGTMImageURL() (string, bool) {
	res, err := http.Get("http://www.lgtm.in/g")
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
		return "not exists", true
	}

	return text, false
}
