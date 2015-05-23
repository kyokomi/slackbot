package suddendeath

import (
	"fmt"
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

func TestCheckMessage(t *testing.T) {
	l := SuddenDeathMessage{}
	ok, _ := l.CheckMessage(context.Background(), ">< 突然の死！")
	if !ok {
		t.Errorf("ERROR check = NG")
	}
}

func TestDoAction(t *testing.T) {
	l := SuddenDeathMessage{}
	ctx := context.Background()
	ctx = plugins.WithSendMessageFunc(ctx, func(message string) {
		fmt.Println(message)
	})
	next := l.DoAction(ctx, ">< 突然の死!")

	if next {
		t.Errorf("ERROR next != false")
	}
}
