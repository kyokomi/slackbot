package router

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kyokomi/slackbot"
	"github.com/kyokomi/slackbot/plugins"
	"github.com/nlopes/slack"
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
	args, ok := ev.BotCmdArgs(message)
	if ok {
		// TODO: Save or Delete or List
		return false
	}
	_ = args

	list, err := p.repository.LoadList(p.repositoryKey(ev.ChannelID()))
	if err != nil {
		ev.Reply(fmt.Sprintf("repository error %s", err.Error()))
		return true
	}

	list = []string{`#random`} // TODO: test

	for _, value := range list {
		channelName, ok := p.filter(value, message)
		if !ok {
			continue
		}

		// TODO: 最初にやってしまってもよいかも?
		chs, err := p.client.GetChannels(true)
		if err != nil {
			ev.Reply(fmt.Sprintf("GetChannelInfo name=[%s] error %s", channelName, err.Error()))
			continue
		}

		for _, ch := range chs {
			// TODO: replacer作っておく
			if ch.Name == strings.Replace(channelName, "#", "", -1) {
				ev.SendMessage(ev.ArchivesURL(), ch.ID)
			}
		}
	}
	return true // next ok
}

func (p *plugin) filter(value string, message string) (string, bool) {
	args := quotationOrSpaceFields(value)
	// hoge.* .*fuga.* #random
	channelName := args[0]

	if len(args) >= 2 {
		expr := args[1]
		if matched, err := regexp.MatchString(expr, message); err != nil || !matched {
			return "", false
		}
	}

	if len(args) >= 3 {
		exclude := args[2]
		if matched, err := regexp.MatchString(exclude, message); err != nil || matched {
			return "", false
		}
	}

	return channelName, true
}

func (p *plugin) repositoryKey(channelID string) string {
	return fmt.Sprintf("router:%s", channelID)
}

func (p *plugin) Help() string {
	return `echo:
	all message echo
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
