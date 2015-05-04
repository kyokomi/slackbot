package lgtm

import (
	"net/http"

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
	return plugins.CheckMessageKeyword(ctx, message, "lgtm")
}

func (m LGTMMessage) DoAction(ctx context.Context, message string) bool {
	sendMessage, isNext := getLGTMImageURL(ctx)

	plugins.SendMessage(ctx, sendMessage)

	return isNext // next stop
}

func getLGTMImageURL(ctx context.Context) (string, bool) {
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

var _ plugins.BotMessagePlugin = (*LGTMMessage)(nil)
