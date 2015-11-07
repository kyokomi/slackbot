package lgtm_test

import (
	"testing"

	"fmt"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/lgtm"
)

var testEvent = plugins.NewTestEvent("LGTM")

func TestCheckMessage(t *testing.T) {
	p := lgtm.NewPlugin()
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := lgtm.NewPlugin()

	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != false")
	}
}

func TestGetLGTMImageURL(t *testing.T) {
	if message, ok := lgtm.GetLGTMImageURL("hoge"); !ok {
		t.Errorf("get lgtm image don't error %s", message)
	} else {
		fmt.Println(message)
	}

	if message, ok := lgtm.GetLGTMImageURL("https://raw.githubusercontent.com/kyokomi/slackbot/master/README.md"); !ok {
		t.Errorf("get lgtm image don't error %s", message)
	} else {
		fmt.Println(message)
	}
}

func TestDoActionArgsUser(t *testing.T) {
	p := lgtm.NewPlugin()

	testEvent := plugins.NewTestEvent("LGTM kyokomi")
	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != false")
	}
}
