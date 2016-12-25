package lgtm

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kyokomi/slackbot/plugins"
)

const lgtmURL = "https://lgtm.in/g"

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

func (p *plugin) Help() string {
	return `lgtm: LGTM
	文中に[LGTM]が含まれていると、LGTM画像をランダムで表示します。

	LGTM <user_name>:
		指定ユーザーのLGTMリストからランダムに表示します。
`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)

func GetLGTMImageURL(lgtmURL string) (string, bool) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{Transport: tr}
	res, err := client.Get(lgtmURL)
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
