package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ActionType string

const (
	AddAction  ActionType = "add"
	DelAction  ActionType = "del"
	StopAction ActionType = "stop"
	ListAction ActionType = "list"
	HelpAction ActionType = "help"
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

	case DelAction, StopAction:
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
	case HelpAction:
		// cron help
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
