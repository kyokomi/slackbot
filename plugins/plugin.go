package plugins

import (
	"fmt"

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
	DoAction(ctx context.Context, message string, sendMessageFunc func(message string)) bool
}

func Context() context.Context {
	return ctx
}

func ExecPlugins(ctx context.Context, message string, sendMessageFunc func(message string)) {
	for _, p := range plugins {
		ok, m := p.CheckMessage(ctx, message)
		if !ok {
			continue
		}

		next := p.DoAction(ctx, m, sendMessageFunc)
		if !next {
			break
		}
	}
}
