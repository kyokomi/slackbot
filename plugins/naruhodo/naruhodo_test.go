package naruhodo_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/naruhodo"
)

var testEvent = plugins.NewTestEvent("それは、なるほど")

func TestCheckMessage(t *testing.T) {
	p := naruhodo.Plugin{}
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := naruhodo.Plugin{}

	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != false")
	}
}
