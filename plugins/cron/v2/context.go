package cron

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/robfig/cron"
)

const helpText = `
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
`

var rd = rand.New(rand.NewSource(time.Now().UnixNano()))

type Context interface {
	AddCommand(channel string, c Command) string
	DelCommand(channel string, c Command) string
	ListCommand(channel string, c Command) string
	HelpCommand(channel string, c Command) string
	Refresh(messageSender plugins.MessageSender, channel string)
	AllRefresh(messageSender plugins.MessageSender)
	Close()
}

type context struct {
	repository  Repository
	cronClient  map[string]*cron.Cron
	cronTaskMap map[string]TaskMap
}

// NewContext create cron context
func NewContext(repository Repository) (Context, error) {
	ctx := &context{
		cronClient: map[string]*cron.Cron{},
		repository: repository,
	}
	data, err := repository.Load()
	if err != nil {
		return nil, err
	}
	ctx.cronTaskMap = data
	return ctx, nil
}

func (ctx *context) AllRefresh(messageSender plugins.MessageSender) {
	for channelID := range ctx.cronTaskMap {
		log.Println("Refresh channelID", channelID)
		ctx.Refresh(messageSender, channelID)
	}
}

func (ctx *context) Close() {
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

func (ctx *context) Refresh(messageSender plugins.MessageSender, channel string) {
	if ctx.cronClient[channel] != nil {
		ctx.cronClient[channel].Stop()
		ctx.cronClient[channel] = nil
	}

	c := cron.New()
	for _, activeCron := range ctx.getTaskMap(channel) {
		if !activeCron.Active {
			continue
		}

		cmd := activeCron.Command
		c.AddFunc(activeCron.Command.CronSpec, func() {
			message := cmd.Message()
			switch cmd.Action {
			case RandomAddAction:
				idx := rd.Intn(len(cmd.Args) - 1)
				log.Println(len(cmd.Args), idx, cmd.Args[idx])
				message = cmd.Args[idx]
			}
			messageSender.SendMessage(message, channel)
		})
	}
	c.Start()

	ctx.cronClient[channel] = c

	if ctx.repository != nil {
		ctx.repository.Save(ctx.cronTaskMap)
	}
}

func (ctx *context) startTask(channelID string, c Command) {
	if ctx.cronTaskMap[channelID] == nil {
		ctx.cronTaskMap[channelID] = TaskMap{}
	}
	ctx.cronTaskMap[channelID].AddTask(c.Key(), Task{true, c})
}

func (ctx *context) stopTask(channelID string, c Command) {
	if ctx.cronTaskMap[channelID] == nil {
		ctx.cronTaskMap[channelID] = TaskMap{}
	}
	ctx.cronTaskMap[channelID].AddTask(c.Key(), Task{false, c})
}

func (ctx *context) getTaskMap(channelID string) map[string]Task {
	if ctx.cronTaskMap[channelID] == nil {
		ctx.cronTaskMap[channelID] = TaskMap{}
	}
	return ctx.cronTaskMap[channelID]
}

func (ctx *context) AddCommand(channel string, c Command) string {
	ctx.startTask(channel, c)
	return fmt.Sprintf("`%s setup done`", c.CronID)
}

func (ctx *context) DelCommand(channel string, c Command) string {
	ctx.stopTask(channel, c)
	return fmt.Sprintf("`%s deleted done`", c.CronID)
}

func (ctx *context) ListCommand(channel string, _ Command) string {
	specMessage := []string{}
	for _, ccd := range ctx.getTaskMap(channel) {
		if !ccd.Active {
			continue
		}
		specMessage = append(specMessage, fmt.Sprintf(
			"cron = [%s] message = [%s] id = [%s]",
			ccd.Command.CronSpec,
			ccd.Command.Message(),
			ccd.Command.CronID,
		))
	}
	message := strings.Join(specMessage, "\n")
	if message == "" {
		message = "not cron list"
	}
	return fmt.Sprintf("```\n%s\n```", message)
}

func (ctx *context) HelpCommand(channel string, _ Command) string {
	return fmt.Sprintf("```\n%s\n```", helpText)
}

type Repository interface {
	Load() (map[string]TaskMap, error)
	Save(cronTaskMap map[string]TaskMap) error
	Close() error
}
