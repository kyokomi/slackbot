package slackbot

import (
	"github.com/kyokomi/slackbot/plugins"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

func messageResponse(ctx context.Context, msEvent *slack.MessageEvent, sendMessageFunc func(message string)) {
	wsAPI := FromSlackRTM(ctx)

	user := wsAPI.GetInfo().User
	if user.Id == msEvent.UserId {
		// 自分のやつはスルーする
		return
	}
	plugins.ExecPlugins(ctx, msEvent, sendMessageFunc)
}
