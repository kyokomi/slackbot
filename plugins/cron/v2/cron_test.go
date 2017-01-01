package cron_test

import (
	"fmt"
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/cron/v2"
)

var testEvents = []plugins.BotEvent{
	plugins.NewTestEvent("cron add */1 * * * * * hogehoge"),
	plugins.NewTestEvent("cron list"),
	plugins.NewTestEvent("cron help"),
	plugins.NewTestEvent("cron del xfjield"),
}

func TestCheckMessage(t *testing.T) {
	repository := cron.NewInMemoryRepository()
	ctx, err := cron.NewContext(repository)
	if err != nil {
		t.Fatal(err.Error())
	}
	p := cron.NewPlugin(ctx)
	for _, testEvent := range testEvents {
		ok, _ := p.CheckMessage(testEvent, testEvent.BaseText())
		if !ok {
			t.Error("ERROR check = NG")
		}
	}
}

func TestDoAction(t *testing.T) {
	repository := cron.NewInMemoryRepository()
	ctx, err := cron.NewContext(repository)
	if err != nil {
		t.Fatal(err.Error())
	}
	p := cron.NewPlugin(ctx)

	for _, testEvent := range testEvents {
		next := p.DoAction(testEvent, testEvent.BaseText())
		if next != false {
			t.Error("ERROR next != false")
		}
	}
}

func TestCronCommand(t *testing.T) {
	command := `cron add */1 * * * * * hogehoge`

	c := cron.Command{}
	if err := c.Scan(command); err != nil {
		t.Errorf("error %s", err)
	}

	if command != fmt.Sprintf("%s", c) {
		t.Errorf("error \n%s\n%s", command, c.String())
	}
}
