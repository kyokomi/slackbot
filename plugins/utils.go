package plugins

import (
	"fmt"
	"strings"
)

// CheckMessageKeyword is Deprecated
func CheckMessageKeyword(message string, keyword string) (bool, string) {
	return CheckMessageKeywords(message, keyword)
}

func CheckMessageKeywords(message string, keywords ...string) (bool, string) {
	for _, keyword := range keywords {
		if nextOK, _ := checkMessageKeyword(message, keyword); nextOK {
			return true, message
		}
	}
	return false, message
}

func checkMessageKeyword(message string, keyword string) (bool, string) {
	if strings.Index(strings.ToLower(message), keyword) == -1 {
		return false, message
	}
	return true, message
}

type DebugMessageSender struct {
}

func (b DebugMessageSender) SendMessage(message string, channel string) {
	fmt.Println(channel)
	fmt.Println(message)
}

var _ MessageSender = (*DebugMessageSender)(nil)

func NewTestEvent(message string) BotEvent {
	return NewBotEvent(DebugMessageSender{},
		"botID", "botName",
		"userID", "userName",
		message,
		"CH_AAAAAAA", "#general",
	)
}
