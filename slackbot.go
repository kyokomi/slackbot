package slackbot

import (
	"errors"
	"fmt"
	"log"

	"github.com/nlopes/slack"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/sysstd"
)

type BotContext interface {
	PluginManager() plugins.PluginManager
	AddPlugin(key interface{}, val plugins.BotMessagePlugin)
	WebSocketRTM()
	Run(connectedFunc func(event plugins.BotEvent))

	// plugins.MessageSender
	SendMessage(message, channel string)
}

type Context struct {
	Plugins plugins.PluginManager
	Client  *slack.Client
	RTM     *slack.RTM
}

func (b *Context) PluginManager() plugins.PluginManager {
	return b.Plugins
}

func NewContext(token string) (*Context, error) {
	ctx, err := NewBotContextNotSysstd(token)
	if err != nil {
		return nil, err
	}
	ctx.AddPlugin("sysstd", sysstd.NewPlugin(ctx.PluginManager()))

	return ctx, nil
}

func NewBotContext(token string) (BotContext, error) {
	ctx, err := NewBotContextNotSysstd(token)
	if err != nil {
		return nil, err
	}
	ctx.AddPlugin("sysstd", sysstd.NewPlugin(ctx.PluginManager()))

	return ctx, nil
}

func NewBotContextNotSysstd(token string) (*Context, error) {
	if token == "" {
		return nil, errors.New("ERROR: slack token not found")
	}

	ctx := &Context{}
	ctx.Client = slack.New(token)
	ctx.Client.SetDebug(true) // TODO: あとで
	ctx.Plugins = plugins.NewPluginManager(ctx)

	return ctx, nil
}

func (ctx *Context) AddPlugin(key interface{}, val plugins.BotMessagePlugin) {
	fmt.Println("insert plugin ", key)
	ctx.Plugins.AddPlugin(key, val)
}

func (ctx *Context) Run(connectedFunc func(event plugins.BotEvent)) {
	ctx.webSocketRTM(connectedFunc)
}

// WebSocketRTM is Deprecated
func (ctx *Context) WebSocketRTM() {
	ctx.webSocketRTM(func(event plugins.BotEvent) { log.Println("connected ", event.Channel()) })
}

func (ctx *Context) webSocketRTM(connectedFunc func(event plugins.BotEvent)) {
	if ctx.RTM != nil {
		ctx.RTM.Disconnect()
	}
	ctx.RTM = ctx.Client.NewRTM()

	go ctx.RTM.ManageConnection()

	go func() {
		for {
			select {
			case msg := <-ctx.RTM.IncomingEvents:
				fmt.Print("Event Received: ")
				switch ev := msg.Data.(type) {
				case *slack.ConnectedEvent:
					botUser := ctx.RTM.GetInfo().User
					for _, c := range ev.Info.Channels {
						connectedFunc(plugins.NewBotEvent(ctx,
							ctx.RTM.GetInfo().Team.Domain,
							botUser.ID, botUser.Name,
							ev.Info.User.ID, ev.Info.User.Name, "connected",
							c.ID, c.Name,
							"",
						))
					}
				case *slack.MessageEvent:
					ctx.Plugins.ExecPlugins(ctx.responseEvent(ev))
				case *slack.PresenceChangeEvent:
					fmt.Printf("Presence Change: %v\n", ev)
				case slack.LatencyReport:
					fmt.Printf("Current latency: %v\n", ev.Value)
				case *slack.RTMError:
					fmt.Printf("Error: %d - %s\n", ev.Code, ev.Msg)
				default:
					fmt.Printf("Unexpected: %+v\n", ev)
				}
			}
		}
	}()
}

func (ctx *Context) responseEvent(ev *slack.MessageEvent) plugins.BotEvent {
	botUser := ctx.RTM.GetInfo().User

	var chName string
	cInfo, err := ctx.RTM.Client.GetChannelInfo(ev.Channel)
	if err == nil {
		chName = cInfo.Name
	}

	return plugins.NewBotEvent(ctx,
		ctx.RTM.GetInfo().Team.Domain,
		botUser.ID,
		botUser.Name,
		ev.User,
		ev.Username,
		ev.Text,
		ev.Channel,
		chName,
		ev.Timestamp,
	)
}

func (ctx *Context) SendMessage(message, channel string) {
	if !ctx.Plugins.IsReply() {
		return
	}
	log.Println("WithSendChannelMessageFunc", channel, message)
	if message != "" {
		ctx.RTM.SendMessage(ctx.RTM.NewOutgoingMessage(message, channel))
	}
}

var _ plugins.MessageSender = (*Context)(nil)
