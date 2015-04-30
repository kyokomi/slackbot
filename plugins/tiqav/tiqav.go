package tiqav

import (
	"fmt"
	"net/http"
	"strings"

	"math/rand"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("tiqav"), TiqavImageMessage{})
}

type TiqavImageMessage struct {
}

func (m TiqavImageMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	if strings.Index(strings.ToLower(message), "image me") != -1 {
		word := strings.Replace(strings.ToLower(message), "image me", "", -1)
		return true, strings.TrimSpace(word)
	}
	return false, message
}

func (m TiqavImageMessage) DoAction(ctx context.Context, message string, sendMessageFunc func(message string)) bool {
	res, err := http.Get("http://api.tiqav.com/search.json?q=" + message)
	if err != nil {
		sendMessageFunc(err.Error())
		return true
	}
	defer res.Body.Close()

	j, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		sendMessageFunc(err.Error())
		return true
	}

	array := j.MustArray()
	if len(array) <= 0 {
		sendMessageFunc("not images")
		return true
	}

	// ランダムに返す
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(array))
	data := array[idx].(map[string]interface{})

	imageURL := fmt.Sprintf("http://img.tiqav.com/%s.%s", data["id"].(string), data["ext"].(string))
	sendMessageFunc(imageURL)
	return false // next stop
}

var _ plugins.BotMessagePlugin = (*TiqavImageMessage)(nil)
