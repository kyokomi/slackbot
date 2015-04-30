package lgtm

import (
	"testing"
	"golang.org/x/net/context"
	"fmt"
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
	next := l.DoAction(context.Background(), "hoge", func(message string) {
		fmt.Println(message)
	})

	if next {
		t.Errorf("ERROR next != false")
	}
}
