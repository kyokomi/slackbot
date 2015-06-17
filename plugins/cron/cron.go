package cron

import (
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
	// TODO: execute message check
	return true, message
}

func (r CronMessage) DoAction(ctx context.Context, message string) bool {
	// TODO: message action
	//plugins.SendMessage(ctx, message)
	return true // next ok
}

var _ plugins.BotMessagePlugin = (*CronMessage)(nil)
