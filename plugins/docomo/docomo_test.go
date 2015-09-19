package docomo_test

import (
	"testing"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/docomo"
)

var testEvent = plugins.NewTestEvent("botID こんにちわ")

func TestCheckMessage(t *testing.T) {
	p := docomo.NewPlugin(nil, slackbot.NewOnMemoryRepository())
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

// TODO: docomo.ClientのMockが必要
//func TestDoAction(t *testing.T) {
//	p := docomo.Plugin{}
//
//	next := p.DoAction(testEvent, testEvent.BaseText())
//
//	if next != true {
//		t.Errorf("ERROR next != true")
//	}
//}
