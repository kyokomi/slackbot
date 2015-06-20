package plugins

import "golang.org/x/net/context"

type key string

const (
	sendMessageKey key = "SendMessageKey"
	sendChannelMessageKey key = "SendChannelMessageKey"
)

func WithSendMessageFunc(ctx context.Context, sendMessageFunc func(message string)) context.Context {
	return context.WithValue(ctx, sendMessageKey, sendMessageFunc)
}

func WithSendChannelMessageFunc(ctx context.Context, sendMessageFunc func(channelID, message string)) context.Context {
	return context.WithValue(ctx, sendChannelMessageKey, sendMessageFunc)
}

func SendMessage(ctx context.Context, message string) {
	if stop {
		return
	}

	sendMessageFunc, ok := ctx.Value(sendMessageKey).(func(message string))
	if ok {
		sendMessageFunc(message)
	}
}

func SendChannelMessage(ctx context.Context, channelID, message string) {
	if stop {
		return
	}

	sendMessageFunc, ok := ctx.Value(sendChannelMessageKey).(func(channelID, message string))
	if ok {
		sendMessageFunc(channelID, message)
	}
}