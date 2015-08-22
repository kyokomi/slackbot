package slackbot

import (
	"fmt"
	"time"

	"log"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/slackctx"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

type BotConfig struct {
	Name       string
	SlackToken string
	Protocol   string
	Origin     string
	KeepAlive  time.Duration
}

func DefaultConfig() BotConfig {
	c := BotConfig{}
	c.Name = "Bot"
	c.Protocol = ""
	c.Origin = "http://example.com"
	c.KeepAlive = 20 * time.Second
	return c
}

func WebSocketRTM(ctx context.Context, config BotConfig) context.Context {
	if config.SlackToken == "" {
		log.Fatal("ERROR: slack token not found")
	}

	ctx = slackctx.NewSlackClient(ctx, config.Name, config.SlackToken)
	ctx = slackctx.NewSlackRTM(ctx, config.Protocol, config.Origin)

	api := slackctx.FromSlackClient(ctx)
	api.SetDebug(true)

	rtm := slackctx.FromSlackRTM(ctx)
	go rtm.ManageConnection()

	ctx = plugins.WithSendChannelMessageFunc(ctx, func(channelID, message string) {
		log.Println("WithSendChannelMessageFunc", channelID, message)
		if message != "" && channelID != "" {
			rtm.SendMessage(rtm.NewOutgoingMessage(message, channelID))
		}
	})

	go func() {
		for {
			select {
			case msg := <-rtm.IncomingEvents:
				fmt.Print("Event Received: ")
				switch ev := msg.Data.(type) {
				case slack.HelloEvent:
				// TODO:
				case *slack.ConnectedEvent:
					fmt.Println("Infos:", ev.Info)
					fmt.Println("Connection counter:", ev.ConnectionCount)
					// Replace #general with your Channel ID
					rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))
				case *slack.MessageEvent:
					messageResponse(ctx, ev, func(message string) {
						if message != "" {
							rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
						}
					})
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

	return ctx
}
