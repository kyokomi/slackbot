package sysstd_test

import (
	"fmt"
	"testing"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/echo"
	"github.com/kyokomi/slackbot/plugins/sysstd"
)

var testEvent = plugins.NewTestEvent("botID date tokyo")

func TestCheckMessage(t *testing.T) {
	p := sysstd.NewPlugin(nil)
	ok, message := p.CheckMessage(testEvent, testEvent.BaseText())
	if !ok {
		t.Errorf("ERROR check = NG")
	} else {
		fmt.Println(message)
	}
}

func TestDoAction(t *testing.T) {
	p := sysstd.NewPlugin(nil)

	next := p.DoAction(testEvent, "date a 1 3")

	if next != false {
		t.Errorf("ERROR next != false")
	}
}

func TestSysstdDebug(t *testing.T) {
	p := sysstd.NewPlugin(nil)
	p.SetDebug(true)
}

func TestSetTimezone(t *testing.T) {
	p := sysstd.NewPlugin(nil)
	p.SetTimezone("JST")
}

func TestSysstdBuildPluginsHelp(t *testing.T) {
	botCtx, err := slackbot.NewBotContext("hoge_token")
	if err != nil {
		t.Error(err)
	}
	botCtx.AddPlugin("echo", echo.NewPlugin())
	p := sysstd.NewPlugin(botCtx.PluginManager())

	ev := plugins.NewTestEvent("botID help")
	p.DoAction(ev, "help")
}
