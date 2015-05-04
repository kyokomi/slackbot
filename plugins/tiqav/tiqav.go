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
	if ok, message := plugins.CheckMessageKeyword(ctx, message, "image me"); ok {
		word := strings.Replace(strings.ToLower(message), "image me", "", -1)
		return true, strings.TrimSpace(word)
	}
	return false, message
}

func (m TiqavImageMessage) DoAction(ctx context.Context, message string) bool {
	imageURL, isNext := getTiqavImageURL(ctx, message)

	plugins.SendMessage(ctx, imageURL)

	return isNext // next stop
}

func getTiqavImageURL(ctx context.Context, message string) (string, bool) {
	res, err := http.Get("http://api.tiqav.com/search.json?q=" + message)
	if err != nil {
		return err.Error(), true
	}
	defer res.Body.Close()

	j, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		return err.Error(), true
	}

	array := j.MustArray()
	if len(array) <= 0 {
		return "not images", true
	}

	// ランダムに返す
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(array))
	data := array[idx].(map[string]interface{})

	imageURL := fmt.Sprintf("http://img.tiqav.com/%s.%s", data["id"].(string), data["ext"].(string))
	return imageURL, false
}

var _ plugins.BotMessagePlugin = (*TiqavImageMessage)(nil)
