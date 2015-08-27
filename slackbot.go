package slackbot

import (
	"errors"
	"fmt"
	"log"

	"github.com/nlopes/slack"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/plugins/sysstd"
)

type BotContext struct {
	Plugins plugins.PluginManager
	Client  *slack.Client
	RTM     *slack.RTM
}

func NewBotContext(token string) (*BotContext, error) {
	if token == "" {
		return nil, errors.New("ERROR: slack token not found")
	}

	ctx := &BotContext{}
	ctx.Client = slack.New(token)
	ctx.Client.SetDebug(true) // TODO: あとで
	ctx.Plugins = plugins.NewPluginManager(ctx)
	ctx.AddPlugin("sysstd", sysstd.Plugin{})

	return ctx, nil
}

func (ctx *BotContext) AddPlugin(key interface{}, val plugins.BotMessagePlugin) {
	fmt.Println("insert plugin ", key)
	ctx.Plugins.AddPlugin(key, val)
}

func (ctx *BotContext) WebSocketRTM() {
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
					fmt.Println("Infos:", ev.Info)
					fmt.Println("Connection counter:", ev.ConnectionCount)
					// Replace #general with your Channel ID
					ctx.SendMessage("Hello world", "#general")
				case *slack.MessageEvent:
					ctx.responseEvent(ev)
				case *slack.PresenceChangeEvent:
					fmt.Printf("Presence Change: %v\n", ev)
				case slack.LatencyReport:
					fmt.Printf("Current latency: %v\n", ev.Value)
				case *slack.RTMError:
					fmt.Printf("Error: %d - %s\n", ev.Code, ev.Msg)
				default:
					fmt.Printf("Unexpected: %+v\n", msg.Data)
				}
			}
		}
	}()
}

func (ctx *BotContext) responseEvent(ev *slack.MessageEvent) {
	botUser := ctx.RTM.GetInfo().User

	e := plugins.NewBotEvent(ctx, botUser.ID, botUser.Name, ev.User, ev.Username, ev.Text, ev.Channel)
	ctx.Plugins.ExecPlugins(e)
}

func (ctx *BotContext) SendMessage(message, channel string) {
	if !ctx.Plugins.IsReply() {
		return
	}
	log.Println("WithSendChannelMessageFunc", channel, message)
	if message != "" && channel != "" {
		ctx.RTM.SendMessage(ctx.RTM.NewOutgoingMessage(message, channel))
	}
}

var _ plugins.MessageSender = (*BotContext)(nil)
