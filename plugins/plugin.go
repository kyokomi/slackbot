package plugins

import (
	"fmt"

	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

var (
	ctx     context.Context
	plugins = []plugin{}
)

type plugin struct {
	Key interface{}
	BotMessagePlugin
}

func AddPlugin(key interface{}, val BotMessagePlugin) {
	fmt.Println("insert plugin ", key)
	plugins = append(plugins, plugin{key, val})
}

func DelPlugin(key interface{}) {
	// TODO: 未実装
	panic("未実装")
}

func init() {
	ctx = context.Background()
}

type BotMessagePlugin interface {
	CheckMessage(ctx context.Context, message string) (bool, string)
	DoAction(ctx context.Context, msEvent *slack.MessageEvent, message string, sendMessageFunc func(message string)) bool
}

func Context() context.Context {
	return ctx
}

func ExecPlugins(ctx context.Context, msEvent *slack.MessageEvent, sendMessageFunc func(message string)) {
	// 条件のfuncとOK時のfunc
	for _, p := range plugins {
		ok, m := p.CheckMessage(ctx, msEvent.Text)
		if !ok {
			continue
		}

		next := p.DoAction(ctx, msEvent, m, sendMessageFunc)
		if !next {
			break
		}
	}
}
