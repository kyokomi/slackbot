package echo_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/echo"
)

var testEvent = plugins.NewTestEvent("てすと")

func TestCheckMessage(t *testing.T) {
	p := echo.NewPlugin()
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := echo.NewPlugin()

	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != true {
		t.Errorf("ERROR next != true")
	}
}
