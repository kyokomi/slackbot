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
	*slack.Client
	Name  string
	Token string
}

func NewSlackClient(ctx context.Context, name string, token string) context.Context {
	c := SlackClient{}
	c.Client = slack.New(token)
	c.Name = name
	c.Token = token
	return context.WithValue(ctx, slackClientKey, c)
}

func FromSlackClient(ctx context.Context) SlackClient {
	return ctx.Value(slackClientKey).(SlackClient)
}

func NewSlackRTM(ctx context.Context, protocol, origin string) context.Context {
	api := FromSlackClient(ctx)
	return context.WithValue(ctx, slackRTMKey, api.NewRTM())
}

func FromSlackRTM(ctx context.Context) *slack.RTM {
	return ctx.Value(slackRTMKey).(*slack.RTM)
}

func WithMessageEvent(ctx context.Context, msEvent *slack.MessageEvent) context.Context {
	return context.WithValue(ctx, messageEventKey, msEvent)
}

func FromMessageEvent(ctx context.Context) *slack.MessageEvent {
	return ctx.Value(messageEventKey).(*slack.MessageEvent)
}
