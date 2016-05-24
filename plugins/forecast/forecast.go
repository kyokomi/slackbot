package forecast

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kyokomi/slackbot/plugins"
)

const endpointURL = "http://weather.livedoor.com/forecast/webservice/json/v1"

//http://weather.livedoor.com/forecast/rss/primary_area.xml
const city = "cityタグ"

type plugin struct {
}

func NewPlugin() plugins.BotMessagePlugin {
	return &plugin{}
}

func (p *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	return plugins.CheckMessageKeyword(message, "天気")
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	values := url.Values{}
	values.Add("city", city)
	response, err := http.Get(endpointURL + "?" + values.Encode())
	if err != nil {
		event.Reply("はて？")
		return false
	}

	defer response.Body.Close()

	body, ioutilErr := ioutil.ReadAll(response.Body)
	if ioutilErr != nil {
		event.Reply("はて？")
		return false
	}

	var str WeatherHacks
	jsonErr := json.Unmarshal(body, &str)
	if jsonErr != nil {
		event.Reply("はて？")
		return false
	}

	resStr := str.Description.Text
	//resStr = strings.Replace(resStr, "。", "にゃ。", -1)	//It's a Joke：）
	event.Reply(resStr)

	return false // next ng
}

func (p *plugin) Help() string {
	return `forecast: 天気予報
天気に反応する
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
