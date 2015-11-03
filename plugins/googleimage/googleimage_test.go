package googleimage_test

import (
	"log"
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/googleimage"
)

var testEvent = plugins.NewTestEvent("image me hoge")

func init() {
	log.SetFlags(log.Llongfile)
}

func TestCheckMessage(t *testing.T) {
	p := googleimage.NewPlugin("cx", "key")
	ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	p := googleimage.NewPlugin("cx", "key")

	next := p.DoAction(testEvent, testEvent.BaseText())

	if next != false {
		t.Errorf("ERROR next != false")
	}
}
