package lgtm

import (
	"fmt"
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

func TestCheckMessage(t *testing.T) {
	l := LGTMMessage{}
	ok, _ := l.CheckMessage(context.Background(), "hoge LGTM desu.")
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	l := LGTMMessage{}
	ctx := context.Background()
	ctx = plugins.WithSendMessageFunc(ctx, func(message string) {
		fmt.Println(message)
	})
	next := l.DoAction(ctx, "hoge")

	if next {
		t.Errorf("ERROR next != false")
	}
}
