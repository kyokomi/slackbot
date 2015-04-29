package slackctx

import (
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

type key string

const (
	slackClientKey  key = "SlackClient"
	slackRTMKey     key = "SlackRTM"
	messageEventKey key = "MessageEventKey"
)

type SlackClient struct {
	*slack.Slack
	Name  string
	Token string
}

func NewSlackClient(ctx context.Context, name string, token string) context.Context {
	c := SlackClient{}
	c.Slack = slack.New(token)
	c.Name = name
	c.Token = token
	return context.WithValue(ctx, slackClientKey, c)
}

func FromSlackClient(ctx context.Context) SlackClient {
	return ctx.Value(slackClientKey).(SlackClient)
}

func NewSlackRTM(ctx context.Context, protocol, origin string) context.Context {
	api := FromSlackClient(ctx)
	wsAPI, err := api.StartRTM(protocol, origin)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, slackRTMKey, wsAPI)
}

func FromSlackRTM(ctx context.Context) *slack.SlackWS {
	return ctx.Value(slackRTMKey).(*slack.SlackWS)
}

func WithMessageEvent(ctx context.Context, msEvent *slack.MessageEvent) context.Context {
	return context.WithValue(ctx, messageEventKey, msEvent)
}

func FromMessageEvent(ctx context.Context) *slack.MessageEvent {
	return ctx.Value(messageEventKey).(*slack.MessageEvent)
}
