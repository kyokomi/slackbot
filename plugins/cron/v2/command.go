package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ActionType string

const (
	AddAction       ActionType = "add"
	RandomAddAction ActionType = "random"
	DelAction       ActionType = "del"
	DeleteAction    ActionType = "delete"
	StopAction      ActionType = "stop"
	ListAction      ActionType = "list"
	RefreshAction   ActionType = "refresh"
	HelpAction      ActionType = "help"
)

type Command struct {
	Trigger  string // cron
	Action   ActionType
	CronSpec string
	Args     []string
	CronID   string
}

func (c Command) Message() string {
	return strings.Join(c.Args, " ")
}

func generateCronID() string {
	return strconv.FormatUint(uint64(time.Now().Unix()), 36)
}

// Scan cron add */1 * * * * * hogehoge -> cronCommand
func (c *Command) Scan(command string) error {
	commands := strings.Fields(command)
	fmt.Println(len(commands), commands)
	c.Trigger = commands[0] // cron
	c.Action = ActionType(commands[1])

	switch c.Action {
	case AddAction, RandomAddAction:
		// cron add 1 * * * * * hogehoge
		if len(commands) < 9 {
			return fmt.Errorf("commands length error %d", len(commands))
		}
		c.CronSpec = strings.Join(commands[2:8], " ")
		c.Args = commands[8:]
		c.CronID = generateCronID()
	case DelAction, DeleteAction, StopAction:
		// cron del <cron_id>
		if len(commands) != 3 {
			return fmt.Errorf("commands length error %d", len(commands))
		}
		c.CronID = commands[2]
	case ListAction, RefreshAction, HelpAction:
		// cron list
		if len(commands) != 2 {
			return fmt.Errorf("commands length error %d", len(commands))
		}
	}

	return nil
}

func (c Command) String() string {
	return fmt.Sprintf("%s %s %s %s", c.Trigger, c.Action, c.CronSpec, c.Message())
}

func (c Command) Key() string {
	return c.CronID
}
