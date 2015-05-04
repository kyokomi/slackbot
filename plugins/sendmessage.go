package plugins

import "golang.org/x/net/context"

type key string

const (
	sendMessageKey key = "SendMessageKey"
)

func WithSendMessageFunc(ctx context.Context, sendMessageFunc func(message string)) context.Context {
	return context.WithValue(ctx, sendMessageKey, sendMessageFunc)
}

func SendMessage(ctx context.Context, message string) {
	sendMessageFunc, ok := ctx.Value(sendMessageKey).(func(message string))
	if ok {
		sendMessageFunc(message)
	}
}
