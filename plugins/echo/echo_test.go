package echo_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/echo"
)

var testEvent = plugins.NewBotEvent(plugins.DebugMessageSender{},
	"bot",
	"user",
	"それは、なるほど。",
	"#general",
)

func TestCheckMessage(t *testing.T) {
	p := echo.Plugin{}
	ok, _ := p.CheckMessage(*testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := echo.Plugin{}

	next := p.DoAction(*testEvent, testEvent.BaseText())

	if next != true {
		t.Errorf("ERROR next != true")
	}
}
