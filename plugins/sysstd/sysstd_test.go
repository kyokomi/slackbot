package sysstd_test

import (
	"fmt"
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/sysstd"
)

var testEvent = plugins.NewTestEvent("botID date tokyo")

func TestCheckMessage(t *testing.T) {
	p := sysstd.NewPlugin()
	ok, message := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	} else {
		fmt.Println(message)
	}
}

func TestDoAction(t *testing.T) {
	p := sysstd.NewPlugin()

	next := p.DoAction(testEvent, "date a 1 3")

	if next != false {
		t.Errorf("ERROR next != false")
	}
}

func TestSysstdDebug(t *testing.T) {
	p := sysstd.NewPlugin()
	p.SetDebug(true)
}

func TestSetTimezone(t *testing.T) {
	p := sysstd.NewPlugin()
	p.SetTimezone("JST")
}
