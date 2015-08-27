package cron

import (
	"log"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
	*CronContext
}

func (p Plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	// cron [action] [cron] [message]
	if strings.HasPrefix(message, "cron") {
		return true, message
	}
	return false, message
}

func (p Plugin) DoAction(event plugins.BotEvent, message string) bool {
	c := CronCommand{}
	if err := c.Scan(message); err != nil {
		log.Printf("error %s", err)
		return false
	}

	switch c.Action {
	case AddAction, RandomAddAction:
		message := p.addCronCommand(event.Channel(), c)
		p.refreshCron(&event, event.Channel())
		event.Reply(message)
	case DelAction, DeleteAction, StopAction:
		message := p.delCronCommand(event.Channel(), c)
		p.refreshCron(&event, event.Channel())
		event.Reply(message)
	case ListAction:
		message := p.listCronCommand(event.Channel(), c)
		event.Reply(message)
	case RefreshAction:
		p.refreshCron(&event, event.Channel())
		event.Reply("```\nrefresh ok\n```")
	case HelpAction:
		message := p.helpCronCommand(event.Channel(), c)
		event.Reply(message)
	}
	return false
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)
