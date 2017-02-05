package plugins

import (
	"fmt"
	"strings"
	"unicode"
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
		"example-domain",
		"botID", "botName",
		"userID", "userName",
		message,
		"CH_AAAAAAA", "#general",
		"1485072037.000010",
	)
}

// Utils util
type Utils struct {
	replacer *strings.Replacer
}

// NewUtils return a utils
func NewUtils(replacer *strings.Replacer) *Utils {
	if replacer == nil {
		replacer = strings.NewReplacer("'", "", `"`, "")
	}
	return &Utils{
		replacer: replacer,
	}
}

// DefaultUtils default utils
var DefaultUtils = NewUtils(nil)

// QuotationOrSpaceFields シングルクォーテーションとダブルクォーテーションを考慮して文字列を空白区切りします
func (u Utils) QuotationOrSpaceFields(s string) []string {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	args := strings.FieldsFunc(s, f)
	for i := range args {
		args[i] = u.replacer.Replace(args[i])
	}
	return args
}
