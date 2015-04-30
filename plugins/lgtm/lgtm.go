package lgtm

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("lgtm"), LGTMMessage{})
}

type LGTMMessage struct {
}

func (m LGTMMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	if strings.Index(strings.ToLower(message), "lgtm") == -1 {
		return false, message
	}
	return true, message
}

func (m LGTMMessage) DoAction(ctx context.Context, message string, sendMessageFunc func(message string)) bool {
	res, err := http.Get("http://www.lgtm.in/g")
	if err != nil {
		sendMessageFunc(err.Error())
		return true
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		sendMessageFunc(err.Error())
		return true
	}
	text, exists := doc.Find("#imageUrl").Attr("value")
	if !exists {
		sendMessageFunc("not exists")
		return true
	}
	sendMessageFunc(text)

	return false // next stop
}

var _ plugins.BotMessagePlugin = (*LGTMMessage)(nil)
