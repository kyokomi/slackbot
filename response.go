package slackbot

import (
	"github.com/kyokomi/slackbot/plugins"
	"github.com/kyokomi/slackbot/slackctx"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

func messageResponse(ctx context.Context, msEvent *slack.MessageEvent, sendMessageFunc func(message string)) {
	wsAPI := slackctx.FromSlackRTM(ctx)

	user := wsAPI.GetInfo().User
	if user.ID == msEvent.BotID {
		// 自分のやつはスルーする
		return
	}

	ctx = slackctx.WithMessageEvent(ctx, msEvent)
	ctx = plugins.WithSendMessageFunc(ctx, sendMessageFunc)

	plugins.ExecPlugins(ctx, msEvent.Text)
}
