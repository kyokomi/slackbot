package sysstd

import (
	"time"
)

const (
	defaultTimeFormat = time.RFC3339
)

var (
	timeZone = map[string]time.Location{
		"JST":   *time.FixedZone("Asia/Tokyo", 9*60*60),
		"UTC":   *time.UTC,
		"Local": *time.Local,
	}
)

type commands []string

var (
	dateCommand        commands = []string{"dateCommand", "date"}
	setTimezoneCommand commands = []string{"setTimezoneCommand", "timezone"}
	execCommand        commands = []string{"command", "cmd", "exec"}
	commandList                 = map[string]commands{
		dateCommand[0]:        dateCommand,
		setTimezoneCommand[0]: setTimezoneCommand,
		execCommand[0]:        execCommand,
	}
)

func (c commands) Contains(text string) bool {
	for _, command := range c {
		if command == text {
			return true
		}
	}
	return false
}
