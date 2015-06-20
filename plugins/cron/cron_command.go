package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

type ActionType string

const (
	AddAction  ActionType = "add"
	DelAction  ActionType = "del"
	ListAction ActionType = "list"
)

type CronCommand struct {
	Trigger  string // cron
	Action   ActionType
	CronSpec string
	Message  string
	CronID   string
}

func generateCronID() string {
	return strconv.FormatUint(uint64(time.Now().Unix()), 36)
}

// Scan cron add */1 * * * * * hogehoge -> cronCommand
func (c *CronCommand) Scan(command string) error {
	commands := strings.Fields(command)
	fmt.Println(len(commands), commands)
	c.Trigger = commands[0] // cron
	c.Action = ActionType(commands[1])

	switch c.Action {
	case AddAction:
		// cron add 1 * * * * * hogehoge
		if len(commands) != 9 {
			return fmt.Errorf("commands length error %d", len(commands))
		}

		c.CronSpec = strings.Join(commands[2:8], " ")
		c.Message = commands[8]
		c.CronID = generateCronID()

	case DelAction:
		// cron del <cron_id>
		if len(commands) != 3 {
			return fmt.Errorf("commands length error %d", len(commands))
		}
		c.CronID = commands[2]

	case ListAction:
		// cron list
		if len(commands) != 2 {
			return fmt.Errorf("commands length error %d", len(commands))
		}
	}

	return nil
}

func (c CronCommand) String() string {
	return fmt.Sprintf("%s %s %s %s", c.Trigger, c.Action, c.CronSpec, c.Message)
}

func (c CronCommand) CronKey() string {
	return c.CronID
}

func addCronCommand(ctx context.Context, channelID string, c CronCommand) {
	startTask(channelID, c)
	refreshCron(ctx, channelID)
	plugins.SendChannelMessage(ctx, channelID, fmt.Sprintf("`%s setup done`", c.CronID))
}

func delCronCommand(ctx context.Context, channelID string, c CronCommand) {
	stopTask(channelID, c)
	refreshCron(ctx, channelID)
	plugins.SendChannelMessage(ctx, channelID, fmt.Sprintf("`%s deleted done`", c.CronID))
}

func listCronCommand(ctx context.Context, channelID string) {
	cronSpecMessage := []string{}
	for _, ccd := range getTaskMap(channelID) {
		if !ccd.Active {
			continue
		}
		cronSpecMessage = append(cronSpecMessage, fmt.Sprintf("%s message = %s id = %s", ccd.Command.CronSpec, ccd.Command.Message, ccd.Command.CronID))
	}
	message := strings.Join(cronSpecMessage, "\n")
	if message == "" {
		message = "not cron list"
	}
	plugins.SendChannelMessage(ctx, channelID, fmt.Sprintf("```\n%s\n```", message))
}
