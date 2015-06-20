package cron

import (
	"log"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
	"github.com/kyokomi/slackbot/slackctx"
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
	msEvent := slackctx.FromMessageEvent(ctx)

	c := CronCommand{}
	if err := c.Scan(message); err != nil {
		log.Printf("error %s", err)
		return false
	}

	switch c.Action {
	case AddAction:
		addCronCommand(ctx, msEvent.ChannelId, c)
	case DelAction:
		delCronCommand(ctx, msEvent.ChannelId, c)
	case ListAction:
		listCronCommand(ctx, msEvent.ChannelId)
	}

	return false
}

var _ plugins.BotMessagePlugin = (*CronMessage)(nil)
