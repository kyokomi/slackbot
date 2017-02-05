package router

import (
	"fmt"
	"strings"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins"
	"github.com/nlopes/slack"
)

const (
	triggerWord = "router"

	helpText = `
	register:
		router add <channel_name> <expr> (<exclude>)
	response:
		<router_id>

	delete:
		router del <router_id>
	response:
		delete routing.

	list:
		router list
	response:
		show added routing list.

	help:
		router help
	response:
		show this help.
`
)

type plugin struct {
	client     *slack.Client
	repository slackbot.Repository
}

func NewPlugin(client *slack.Client, repository slackbot.Repository) plugins.BotMessagePlugin {
	return &plugin{
		client:     client,
		repository: repository,
	}
}

func (p *plugin) CheckMessage(ev plugins.BotEvent, message string) (bool, string) {
	return true, message
}

func (p *plugin) DoAction(ev plugins.BotEvent, message string) bool {
	c := Command{}
	ok, err := c.Scan(message)
	if err != nil {
		ev.Reply("command scan error " + err.Error())
		return false
	}

	if ok && c.Trigger == triggerWord {
		switch c.Action {
		case AddAction:
			addFilter := newFilter(c.RouterID, c.FilterArgs)
			if err := p.saveAddFilter(ev.ChannelID(), addFilter); err != nil {
				ev.Reply("AddAction error " + err.Error())
				return false
			}
			ev.Reply("save rotuter " + c.Message())
		case DelAction, DeleteAction:
			if err := p.saveDeleteFilter(ev.ChannelID(), c.RouterID); err != nil {
				ev.Reply("DelAction error " + err.Error())
				return false
			}
			ev.Reply("delete router " + c.RouterID)
		case ListAction:
			fs, err := p.loadFilters(ev.ChannelID())
			if err != nil {
				ev.Reply("ListAction error " + err.Error())
				return false
			}

			if len(fs) > 0 {
				message := []string{}
				for _, f := range fs {
					message = append(message, f.String())
				}
				ev.Reply(strings.Join(message, "\n"))
			} else {
				ev.Reply("not routing")
			}
		case HelpAction:
			ev.Reply(fmt.Sprintf("```\n%s\n```", helpText))
		default:
			ev.Reply(fmt.Sprintf("not support command [%s]", c.String()))
		}
		return false
	}

	fs, err := p.loadFilters(ev.ChannelID())
	if err != nil {
		// 通常のメッセージが通るのでスルーする
		return true // next ok
	}

	for _, f := range fs {
		if ev.SenderID() == ev.BotID() {
			continue
		}
		if f.ChannelName == ev.ChannelName() {
			continue
		}
		if ok := f.execute(message); !ok {
			continue
		}

		chs, err := p.client.GetChannels(true)
		if err != nil {
			ev.Reply(fmt.Sprintf("GetChannelInfo name=[%s] error %s", f.ChannelName, err.Error()))
			continue
		}

		for _, ch := range chs {
			if ch.Name == strings.Replace(f.ChannelName, "#", "", -1) {
				ev.SendMessage(ev.ArchivesURL(), ch.ID)
			}
		}
	}
	return true // next ok
}

func (p *plugin) repositoryKey(channelID string) string {
	return fmt.Sprintf("router:%s", channelID)
}

func (p *plugin) Help() string {
	return `router:
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
