package tiqav

import (
	"fmt"
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"golang.org/x/net/context"
)

func TestCheckMessage(t *testing.T) {
	l := TiqavImageMessage{}
	ok, message := l.CheckMessage(context.Background(), "image me Gopher")
	if !ok {
		t.Errorf("ERROR check = NG")
	}

	if message != "gopher" {
		t.Errorf("ERROR message %s != %s", message, "gopher")
	}
}

func TestDoAction(t *testing.T) {
	l := TiqavImageMessage{}
	ctx := context.Background()
	plugins.WithSendMessageFunc(ctx, func(message string) {
		fmt.Println(message)
	})
	next := l.DoAction(ctx, "gopher")

	if next {
		t.Errorf("ERROR next != false")
	}
}
