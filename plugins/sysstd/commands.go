package sysstd

import (
	"time"
)

const (
	defaultTimeFormat = time.RFC3339
)

type commandFunc func(args ...string) string

type commands struct {
	commandList []string
	commandFunc commandFunc
}

func (c commands) Contains(text string) bool {
	for _, command := range c.commandList {
		if command == text {
			return true
		}
	}
	return false
}
