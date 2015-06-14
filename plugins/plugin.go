package plugins

import (
	"fmt"

	"golang.org/x/net/context"
)

var (
	ctx     context.Context
	plugins = []Plugin{}
	stop    = false
)

type Plugin struct {
	Key interface{}
	BotMessagePlugin
}

func (p Plugin) Name() string {
	return fmt.Sprintf("%s", p.Key)
}

func AddPlugin(key interface{}, val BotMessagePlugin) {
	fmt.Println("insert plugin ", key)
	plugins = append(plugins, Plugin{key, val})
}

func DelPlugin(key interface{}) {
	// TODO: 未実装
	panic("未実装")
}

// Stop bot stop
func Stop() {
	stop = true
}

// Start bot start
func Start() {
	stop = false
}

// GetPlugins returns all plugins
func GetPlugins() []Plugin {
	return plugins
}

func init() {
	ctx = context.Background()
}

type BotMessagePlugin interface {
	CheckMessage(ctx context.Context, message string) (bool, string)
	DoAction(ctx context.Context, message string) bool
}

func Context() context.Context {
	return ctx
}

func ExecPlugins(ctx context.Context, message string) {
	if stop {
		return
	}

	for _, p := range plugins {
		ok, m := p.CheckMessage(ctx, message)
		if !ok {
			continue
		}

		next := p.DoAction(ctx, m)
		if !next {
			break
		}
	}
}
