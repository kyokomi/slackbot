package cron

import (
	"log"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type pluginKey string

func init() {
	plugins.AddPlugin(pluginKey("cronMessage"), CronMessage{})
}

type CronMessage struct {
}

func (r CronMessage) CheckMessage(ctx context.Context, message string) (bool, string) {
	// cron [action] [cron] [message]
	if strings.HasPrefix(message, "cron") {
		return true, message
	}
	return false, message
}

func (r CronMessage) DoAction(ctx context.Context, message string) bool {
	if _, ok := cronTaskMap[message]; ok {
		return false
	}

	c := cronCommand{}
	if err := c.Scan(message); err != nil {
		log.Printf("error %s", err)
		return false
	}

	switch c.Action {
	case AddAction:
		addCronCommand(ctx, c)
	case DelAction:
		delCronCommand(ctx, c)
	case ListAction:
		listCronCommand(ctx)
	}

	return false
}

var _ plugins.BotMessagePlugin = (*CronMessage)(nil)
