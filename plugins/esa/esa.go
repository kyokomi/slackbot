package esa

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/upamune/go-esa/esa"
)

var urlMarksReplacer = strings.NewReplacer("<", "", ">", "")

type plugin struct {
	teamName  string
	esaClient *esa.Client
}

// NewPlugin esa plugin
func NewPlugin(teamName, token string) plugins.BotMessagePlugin {
	return &plugin{
		teamName:  teamName,
		esaClient: esa.NewClient(token),
	}
}

// CheckMessage esa.io url is ok
func (p *plugin) CheckMessage(_ plugins.BotEvent, message string) (bool, string) {
	// esaのURLが見つかったら1件目を返す
	fields := plugins.DefaultUtils.QuotationOrSpaceFields(message)
	for _, val := range fields {
		u, err := url.Parse(urlMarksReplacer.Replace(val))
		if err != nil || !strings.HasSuffix(u.Host, "esa.io") {
			continue
		}
		return true, u.String()
	}
	return false, message
}

// DoAction is replay url detail message
func (p *plugin) DoAction(event plugins.BotEvent, message string) bool {
	u, _ := url.Parse(message)

	var postNumber int
	paths := strings.Split(u.Path, "/")
	for i := range paths {
		if paths[i] == "posts" && len(paths) >= i+1 {
			postNumber, _ = strconv.Atoi(paths[i+1])
			break
		}
	}
	resp, err := p.esaClient.Post.GetPost(p.teamName, postNumber)
	if err != nil {
		event.Reply(fmt.Sprintf("GetPost error %s", err.Error()))
		return true
	}
	event.Reply(fmt.Sprintf("```\n%s\n```", resp.FullName+"\n"+resp.CreatedAt+"\n\n"+resp.BodyMd))
	return true // next ok
}

// Help print help
func (p *plugin) Help() string {
	return `esa:
	URLを貼ると詳細を展開します
	`
}

var _ plugins.BotMessagePlugin = (*plugin)(nil)
