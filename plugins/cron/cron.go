package cron

import (
	"fmt"
	"log"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/robfig/cron"
	"golang.org/x/net/context"
)

type pluginKey string

var cronClient *cron.Cron

var cronTaskMap = map[string]cronTask{}

type cronTask struct {
	active  bool
	command cronCommand
}

func init() {
	plugins.AddPlugin(pluginKey("cronMessage"), CronMessage{})
}

func Setup() {
	cronClient = cron.New()
	cronClient.Start()
}

func Stop() {
	if cronClient != nil {
		cronClient.Stop()
	}
}

func addCronCommand(ctx context.Context, c cronCommand) {
	cronTaskMap[c.CronKey()] = cronTask{true, c}
	refreshCron(ctx)
	plugins.SendMessage(ctx, fmt.Sprintf("%s setup done", c))
}

func delCronCommand(ctx context.Context, c cronCommand) {
	cronTaskMap[c.CronKey()] = cronTask{false, c}
	refreshCron(ctx)
	plugins.SendMessage(ctx, fmt.Sprintf("%s deleted done", c))
}

func listCronCommand(ctx context.Context) {
	cronSpecMessage := []string{}
	for _, ccd := range cronTaskMap {
		if !ccd.active {
			continue
		}
		cronSpecMessage = append(cronSpecMessage, fmt.Sprintf("%s : %s", ccd.command.CronSpec, ccd.command.Message))
	}
	plugins.SendMessage(ctx, strings.Join(cronSpecMessage, "\n"))
}

func refreshCron(ctx context.Context) {
	cronClient.Stop()
	c := cron.New()
	for _, activeCron := range cronTaskMap {
		if !activeCron.active {
			continue
		}

		message := activeCron.command.Message
		c.AddFunc(activeCron.command.CronSpec, func() {
			plugins.SendMessage(ctx, message)
		})
	}
	c.Start()
	cronClient = c
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

type ActionType string

const (
	AddAction  ActionType = "add"
	DelAction  ActionType = "del"
	ListAction ActionType = "list"
)

type cronCommand struct {
	Trigger  string // cron
	Action   ActionType
	CronSpec string
	Message  string
}

// Scan cron add */1 * * * * * hogehoge -> cronCommand
func (c *cronCommand) Scan(command string) error {
	commands := strings.Fields(command)
	fmt.Println(len(commands), commands)
	c.Trigger = commands[0]
	c.Action = ActionType(commands[1])

	if c.Action == AddAction || c.Action == DelAction {
		if len(commands) != 9 {
			return fmt.Errorf("commands length error %d", len(commands))
		}

		c.CronSpec = strings.Join(commands[2:8], " ")
		c.Message = commands[8]
	} else if c.Action == ListAction {
		if len(commands) != 2 {
			return fmt.Errorf("commands length error %d", len(commands))
		}
	}

	return nil
}

func (c cronCommand) String() string {
	return fmt.Sprintf("%s %s %s %s", c.Trigger, c.Action, c.CronSpec, c.Message)
}

func (c cronCommand) CronKey() string {
	return fmt.Sprintf("%s %s", c.CronSpec, c.Message)
}
