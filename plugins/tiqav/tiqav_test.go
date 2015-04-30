package tiqav

import (
	"testing"
	"golang.org/x/net/context"
	"fmt"
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
	next := l.DoAction(context.Background(), "gopher", func(message string) {
		fmt.Println(message)
	})

	if next {
		t.Errorf("ERROR next != false")
	}
}
