package router

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ActionType string

const (
	AddAction    ActionType = "add"
	DelAction    ActionType = "del"
	DeleteAction ActionType = "delete"
	ListAction   ActionType = "list"
	HelpAction   ActionType = "help"
)

type Command struct {
	Trigger    string // router
	Action     ActionType
	FilterArgs string
	RouterID   string
}

func generateRouterID() string {
	return strconv.FormatUint(uint64(time.Now().Unix()), 36)
}

// Scan router add hoge.* .*fuga.* #random -> cronCommand
func (c *Command) Scan(command string) (bool, error) {
	commands := strings.Fields(command)
	if len(commands) < 2 {
		return false, nil
	}
	c.Trigger = commands[0]
	c.Action = ActionType(commands[1])

	switch c.Action {
	case AddAction:
		// router add hoge.* .*fuga.* #random
		if len(commands) < 4 || 5 < len(commands) {
			return false, fmt.Errorf("commands length error %d", len(commands))
		}
		c.FilterArgs = strings.Join(commands[2:], " ")
		c.RouterID = generateRouterID()
	case DelAction, DeleteAction:
		// router del xxxxxxxxx
		if len(commands) != 3 {
			return false, fmt.Errorf("commands length error %d", len(commands))
		}
		c.RouterID = commands[2]
	case ListAction, HelpAction:
		// router list
		if len(commands) != 2 {
			return false, fmt.Errorf("commands length error %d", len(commands))
		}
	}

	return true, nil
}

func (c Command) String() string {
	return fmt.Sprintf("%s %s %s %s", c.Trigger, c.Action, c.FilterArgs, c.RouterID)
}

// Message return a reply message
func (c Command) Message() string {
	return fmt.Sprintf("FilterArgs = [%s] ID = [%s]", c.FilterArgs, c.RouterID)
}

func (c Command) Key() string {
	return c.RouterID
}
