package lgtm_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/lgtm"
)

var testEvent = plugins.NewBotEvent(plugins.DebugMessageSender{},
	"bot",
	"user",
	"LGTM",
	"#general",
)

func TestCheckMessage(t *testing.T) {
	p := lgtm.Plugin{}
	ok, _ := p.CheckMessage(*testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := lgtm.Plugin{}

	next := p.DoAction(*testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != false")
	}
}
