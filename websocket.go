package slackbot

import (
	"fmt"
	"time"

	"log"

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

func WebSocketRTM(ctx context.Context, config BotConfig) {
	if config.SlackToken == "" {
		log.Fatal("ERROR: slack token not found")
	}

	ctx = NewSlackClient(ctx, config.Name, config.SlackToken)
	ctx = NewSlackRTM(ctx, config.Protocol, config.Origin)

	chSender := make(chan slack.OutgoingMessage)
	chReceiver := make(chan slack.SlackEvent)

	api := FromSlackClient(ctx)
	api.SetDebug(true)

	wsAPI := FromSlackRTM(ctx)
	go wsAPI.HandleIncomingEvents(chReceiver)
	go wsAPI.Keepalive(config.KeepAlive)
	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, chSender)

	go func(chSender chan slack.OutgoingMessage, chReceiver chan slack.SlackEvent) {
		for {
			select {
			case msg := <-chReceiver:
				fmt.Print("Event Received: ")
				switch msg.Data.(type) {
				case slack.HelloEvent:
				// TODO:
				case *slack.MessageEvent:
					a := msg.Data.(*slack.MessageEvent)
					messageResponse(ctx, a, func(message string) {
						if message != "" {
							chSender <- *wsAPI.NewOutgoingMessage(message, a.ChannelId)
						}
					})
				case *slack.PresenceChangeEvent:
					a := msg.Data.(*slack.PresenceChangeEvent)
					fmt.Printf("Presence Change: %v\n", a)
				case slack.LatencyReport:
					a := msg.Data.(slack.LatencyReport)
					fmt.Printf("Current latency: %v\n", a.Value)
				case *slack.SlackWSError:
					error := msg.Data.(*slack.SlackWSError)
					fmt.Printf("Error: %d - %s\n", error.Code, error.Msg)
				default:
					fmt.Printf("Unexpected: %v\n", msg.Data)
				}
			}
		}
	}(chSender, chReceiver)
}
