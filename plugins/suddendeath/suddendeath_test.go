package suddendeath_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/suddendeath"
)

var testEvent = plugins.NewTestEvent("突然の死だああああああああああ！")

func TestCheckMessage(t *testing.T) {
	p := suddendeath.NewPlugin()
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := suddendeath.NewPlugin()

	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != false")
	}
}
