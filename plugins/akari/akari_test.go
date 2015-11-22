package akari_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/akari"
)

var testEvent = plugins.NewTestEvent("test大好き")

func TestCheckMessage(t *testing.T) {
	p := akari.NewPlugin()
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := akari.NewPlugin()

	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != false")
	}
}
