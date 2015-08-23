package suddendeath_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/suddendeath"
)

var testEvent = plugins.NewBotEvent(plugins.DebugMessageSender{},
	"bot",
	"user",
	"突然の死だああああああああああ！",
	"#general",
)

func TestCheckMessage(t *testing.T) {
	p := suddendeath.Plugin{}
	ok, _ := p.CheckMessage(*testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := suddendeath.Plugin{}

	next := p.DoAction(*testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != false")
	}
}
