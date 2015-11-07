package cron

import (
	"log"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
)

type plugin struct {
	cron CronContext
}

func NewPlugin(cron CronContext) plugins.BotMessagePlugin {
	return &plugin{
		cron: cron,
	}
}

func (p *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	// cron [action] [cron] [message]
	if strings.HasPrefix(message, "cron") {
		return true, message
	}
	return false, message
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	c := CronCommand{}
	if err := c.Scan(message); err != nil {
		log.Printf("error %s", err)
		return false
	}

	switch c.Action {
	case AddAction, RandomAddAction:
		message := p.cron.AddCronCommand(event.Channel(), c)
		p.cron.RefreshCron(&event, event.Channel())
		event.Reply(message)
	case DelAction, DeleteAction, StopAction:
		message := p.cron.DelCronCommand(event.Channel(), c)
		p.cron.RefreshCron(&event, event.Channel())
		event.Reply(message)
	case ListAction:
		message := p.cron.ListCronCommand(event.Channel(), c)
		event.Reply(message)
	case RefreshAction:
		p.cron.RefreshCron(&event, event.Channel())
		event.Reply("```\nrefresh ok\n```")
	case HelpAction:
		message := p.cron.HelpCronCommand(event.Channel(), c)
		event.Reply(message)
	}
	return false
}

func (p *plugin) Help() string {
	return "cron: cron制御できます\n" + helpText
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
