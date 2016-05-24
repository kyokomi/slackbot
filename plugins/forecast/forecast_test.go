package forecast_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/forecast"
)

var testEvent = plugins.NewTestEvent("天気にゃ")

func TestCheckMessage(t *testing.T) {
	p := forecast.NewPlugin()
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}
