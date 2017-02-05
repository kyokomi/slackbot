package esa_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/esa"
)

var testEvent = plugins.NewTestEvent("ほげ http://example.esa.io/posts/111111 これはtests")

func TestCheckMessage(t *testing.T) {
	p := esa.NewPlugin("example", "")
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Error("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := esa.NewPlugin("example", "")

	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != true {
		t.Error("ERROR next != true")
	}
}
