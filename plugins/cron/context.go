package cron

import (
	"fmt"
	"log"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/robfig/cron"
)

type CronContext struct {
	repository  CronRepository
	cronClient  map[string]*cron.Cron
	cronTaskMap map[string]CronTaskMap
}

func NewCronContext(repository CronRepository) *CronContext {
	ctx := &CronContext{
		cronClient: map[string]*cron.Cron{},
		repository: repository,
	}
	data, err := repository.Load()
	if err != nil {
		log.Println(err)
	} else {
		ctx.cronTaskMap = data
	}
	return ctx
}

func (ctx *CronContext) AllRefreshCron(messageSender plugins.MessageSender) {
	for channelID, _ := range ctx.cronTaskMap {
		log.Println("RefreshCron channelID", channelID)
		ctx.refreshCron(messageSender, channelID)
	}
}

func (ctx *CronContext) Close() {
	if ctx.cronClient != nil {
		for _, c := range ctx.cronClient {
			if c == nil {
				continue
			}
			c.Stop()
		}
	}

	if ctx.repository != nil {
		ctx.repository.Close()
	}
}

func (ctx *CronContext) refreshCron(messageSender plugins.MessageSender, channel string) {
	if ctx.cronClient[channel] != nil {
		ctx.cronClient[channel].Stop()
		ctx.cronClient[channel] = nil
	}

	c := cron.New()
	for _, activeCron := range ctx.getTaskMap(channel) {
		if !activeCron.Active {
			continue
		}

		message := activeCron.Command.Message
		c.AddFunc(activeCron.Command.CronSpec, func() {
			messageSender.SendMessage(message, channel)
		})
	}
	c.Start()

	ctx.cronClient[channel] = c

	if ctx.repository != nil {
		ctx.repository.Save(ctx.cronTaskMap)
	}
}

func (ctx *CronContext) startTask(channelID string, c CronCommand) {
	if ctx.cronTaskMap[channelID] == nil {
		ctx.cronTaskMap[channelID] = CronTaskMap{}
	}
	ctx.cronTaskMap[channelID].AddTask(c.CronKey(), CronTask{true, c})
}

func (ctx *CronContext) stopTask(channelID string, c CronCommand) {
	if ctx.cronTaskMap[channelID] == nil {
		ctx.cronTaskMap[channelID] = CronTaskMap{}
	}
	ctx.cronTaskMap[channelID].AddTask(c.CronKey(), CronTask{false, c})
}

func (ctx *CronContext) getTaskMap(channelID string) map[string]CronTask {
	if ctx.cronTaskMap[channelID] == nil {
		ctx.cronTaskMap[channelID] = CronTaskMap{}
	}
	return ctx.cronTaskMap[channelID]
}

func (ctx *CronContext) addCronCommand(channel string, c CronCommand) string {
	ctx.startTask(channel, c)
	return fmt.Sprintf("`%s setup done`", c.CronID)
}

func (ctx *CronContext) delCronCommand(channel string, c CronCommand) string {
	ctx.stopTask(channel, c)
	return fmt.Sprintf("`%s deleted done`", c.CronID)
}

func (ctx *CronContext) listCronCommand(channel string, _ CronCommand) string {
	cronSpecMessage := []string{}
	for _, ccd := range ctx.getTaskMap(channel) {
		if !ccd.Active {
			continue
		}
		cronSpecMessage = append(cronSpecMessage, fmt.Sprintf("%s message = %s id = %s", ccd.Command.CronSpec, ccd.Command.Message, ccd.Command.CronID))
	}
	message := strings.Join(cronSpecMessage, "\n")
	if message == "" {
		message = "not cron list"
	}
	return fmt.Sprintf("```\n%s\n```", message)
}

func (ctx *CronContext) helpCronCommand(channel string, _ CronCommand) string {
	return fmt.Sprintf("```%s```", `
register:
	cron add */1 * * * * * hogehoge
response:
	<cron_id>

delete:
	cron del <cron_id>
response:
	delete message.

list:
	cron list
response:
	show added cron list.

help:
	cron help
response:
	show this help.
`)
}

type CronRepository interface {
	Load() (map[string]CronTaskMap, error)
	Save(cronTaskMap map[string]CronTaskMap) error
	Close() error
}

type CronTask struct {
	Active  bool
	Command CronCommand
}

type CronTaskMap map[string]CronTask

func (c CronTaskMap) AddTask(key string, task CronTask) {
	c[key] = task
}
