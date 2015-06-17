package cron

import (
	"fmt"
	"strings"

	"log"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/robfig/cron"
	"golang.org/x/net/context"
)

type pluginKey string

var cronClient *cron.Cron

var cronTaskMap = map[string]cronTask{}

type cronTask struct {
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
		cronTaskMap[message] = cronTask{c}
		cronClient.AddFunc(c.CronSpec, func() { plugins.SendMessage(ctx, c.Message) })
		plugins.SendMessage(ctx, fmt.Sprintf("%s setup done", c))
	}

	return false
}

var _ plugins.BotMessagePlugin = (*CronMessage)(nil)

type ActionType string

const (
	AddAction ActionType = "add"
)

type cronCommand struct {
	Trigger  string // cron
	Action   ActionType
	CronSpec string
	Message  string
}

// Scan cron add "*/1 * * * * *" hogehoge -> cronCommand
func (c *cronCommand) Scan(command string) error {
	commands := strings.Fields(command)
	fmt.Println(len(commands), commands)
	if len(commands) != 9 {
		return fmt.Errorf("commands length error %d", len(commands))
	}
	c.Trigger = commands[0]
	c.Action = ActionType(commands[1])
	c.CronSpec = strings.Join(commands[2:8], " ")
	c.Message = commands[8]

	return nil
}

func (c cronCommand) String() string {
	return fmt.Sprintf("%s %s %s %s", c.Trigger, c.Action, c.CronSpec, c.Message)
}
