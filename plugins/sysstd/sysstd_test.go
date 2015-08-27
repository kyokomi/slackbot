package sysstd_test

import (
	"testing"
	"fmt"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/sysstd"
)

var testEvent = plugins.NewTestEvent("botID date tokyo")

func TestCheckMessage(t *testing.T) {
	p := sysstd.Plugin{}
	ok, message := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	} else {
		fmt.Println(message)
	}
}

func TestDoAction(t *testing.T) {
	p := sysstd.Plugin{}

	next := p.DoAction(testEvent, "dateCommand a 1 3")

	if next != false {
		t.Errorf("ERROR next != false")
	}
}
