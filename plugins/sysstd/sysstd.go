package sysstd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"bytes"

	"github.com/kyokomi/slackbot/plugins"
)

type Plugin interface {
	SetDebug(debug bool)
	SetTimezone(tmText string)

	// plugins.BotMessagePlugin
	CheckMessage(event plugins.BotEvent, message string) (bool, string)
	DoAction(event plugins.BotEvent, message string) bool
	Help() string
}

type plugin struct {
	pm plugins.PluginManager

	debug       bool
	timeZone    string
	timeZoneMap map[string]time.Location

	commandMap map[string]commands
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
var _ Plugin = (*plugin)(nil)

func NewPlugin(pm plugins.PluginManager) Plugin {
	p := &plugin{
		pm:       pm,
		debug:    false,
		timeZone: "Local",
		timeZoneMap: map[string]time.Location{
			"JST":   *time.FixedZone("Asia/Tokyo", 9*60*60),
			"UTC":   *time.UTC,
			"Local": *time.Local,
		},
	}

	p.commandMap = map[string]commands{
		"date": commands{
			commandList: []string{"date", "dateCommand"},
			commandFunc: func(args ...string) string {
				loc := p.getTimeZoneLocation()
				return time.Now().In(&loc).Format(defaultTimeFormat)
			},
		},
		"help": commands{
			commandList: []string{"help", "-h", "--help"},
			commandFunc: func(args ...string) string {
				return p.buildPluginsHelp()
			},
		},
		"timezone": commands{
			commandList: []string{"timezone", "setTimezoneCommand"},
			commandFunc: func(args ...string) string {
				if len(args) >= 2 {
					p.SetTimezone(args[1])
				}
				return fmt.Sprintf("%#v", p.getTimeZoneLocation())
			},
		},
		"command": commands{
			commandList: []string{"command", "cmd", "exec"},
			commandFunc: func(args ...string) string {
				if len(args) < 3 {
					return "command not found"
				}
				data, err := exec.Command(args[1], args[2:]...).CombinedOutput()
				if err != nil {
					return fmt.Sprintf("`%s`", err.Error())
				} else {
					return fmt.Sprintf("```\n%s```", string(data))
				}
			},
		},
	}

	return p
}

func (p *plugin) buildPluginsHelp() string {
	resMessage := bytes.Buffer{}
	resMessage.WriteString("```\n")
	for idx, p := range p.pm.GetPlugins() {
		if idx != 0 {
			resMessage.WriteString("\n")
		}
		resMessage.WriteString(p.Help())
	}
	resMessage.WriteString("\n```")
	return resMessage.String()
}

func (p *plugin) SetDebug(debug bool) {
	p.debug = debug
}

func (p *plugin) SetTimezone(tmText string) {
	timeZone := strings.ToUpper(tmText)
	if _, ok := p.timeZoneMap[timeZone]; ok {
		p.timeZone = timeZone
	}
}

func (p *plugin) getTimeZoneLocation() time.Location {
	return p.timeZoneMap[p.timeZone]
}

func (p *plugin) executeCommand(args ...string) string {
	for cmdKey, cmds := range p.commandMap {
		if !cmds.Contains(args[0]) {
			continue
		}
		return p.commandMap[cmdKey].commandFunc(args...)
	}
	return "`command error`"
}

func (p *plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	if p.debug {
		log.Printf("message            [%s]\n", message)
		log.Printf("botLinkID          [%s]\n", event.BotLinkID())
		log.Printf("botLinkIDForClient [%s]\n", event.BotLinkIDForClient())
		log.Printf("botName            [%s]\n", event.BotName())
		log.Printf("botID              [%s]\n", event.BotID())
	}

	cmdArgs, ok := event.BotCmdArgs(message)
	if !ok {
		return false, message
	}

	if len(cmdArgs) < 1 {
		return false, message
	}

	if p.debug {
		log.Println(cmdArgs)
	}

	for cmdKey, cmds := range p.commandMap {
		if cmds.Contains(cmdArgs[0]) {
			cmdArgs[0] = cmdKey
			return true, strings.Join(cmdArgs, " ")
		}
	}

	return false, message
}

func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	cmdArgs := strings.Fields(message)
	if _, ok := p.commandMap[cmdArgs[0]]; !ok {
		return true
	}
	event.Reply(p.executeCommand(cmdArgs...))
	return false // next ok
}
func (p *plugin) Help() string {
	return `sysstd: botシステム系のコマンド

	<botName/botID> date :

		現在時間を表示

	<botName/botID> help :

		すべてのプラグインのHelpを表示

	<botName/botID> timezone :

		現在設定されているタイムゾーンを表示

	<botName/botID> timezone <JST/UTC/Local>:

		指定したタイムゾーンに変更

	<botName/botID> command <linux command>:

		bot環境によって動かないコマンドもあります。

		Example:

			botName command ls -al
	`
}
