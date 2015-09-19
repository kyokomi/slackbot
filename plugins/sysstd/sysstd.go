package sysstd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
	debug bool
}

func (p Plugin) SetTimezone(tmText string) {
	reTimeZone, ok := timeZone[strings.ToUpper(tmText)]
	if ok {
		// TODO: かなり強引なので注意...BotContextでTimezone持つべき
		time.Local = &reTimeZone
	}
}

func (p Plugin) ExecuteCommand(args ...string) string {
	switch {
	case execCommand.Contains(args[0]):
		if len(args) < 3 {
			return "command not found"
		}
		data, err := exec.Command(args[1], args[2:]...).CombinedOutput()
		if err != nil {
			return fmt.Sprintf("`%s`", err.Error())
		} else {
			return fmt.Sprintf("```\n%s```", string(data))
		}
	case dateCommand.Contains(args[0]):
		return time.Now().Format(defaultTimeFormat)
	case setTimezoneCommand.Contains(args[0]):
		if len(args) >= 2 {
			p.SetTimezone(args[1])
		}
		return fmt.Sprintf("%#v", *time.Local)
	}
	return "`command error`"
}

func (r Plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	if r.debug {
		log.Printf("message   [%s]\n", message)
		log.Printf("botLinkID [%s]\n", event.BotLinkID())
		log.Printf("botName   [%s]\n", event.BotName())
		log.Printf("botID     [%s]\n", event.BotID())
	}

	cmdArgs, ok := event.BotCmdArgs(message)
	if !ok {
		return false, message
	}

	if len(cmdArgs) < 1 {
		return false, message
	}

	if r.debug {
		log.Println(cmdArgs)
	}

	for cmdKey, cmd := range commandList {
		if cmd.Contains(cmdArgs[0]) {
			cmdArgs[0] = cmdKey
			return true, strings.Join(cmdArgs, " ")
		}
	}

	return false, message
}

func (r Plugin) DoAction(event plugins.BotEvent, message string) bool {
	cmdArgs := strings.Fields(message)
	if _, ok := commandList[cmdArgs[0]]; !ok {
		return true
	}
	event.Reply(r.ExecuteCommand(cmdArgs...))
	return false // next ok
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)
