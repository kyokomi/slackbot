package router

import (
	"fmt"
	"regexp"

	"github.com/kyokomi/slackbot/plugins"
)

type filter struct {
	ID          string `json:"id"`
	ChannelName string `json:"channelName"`
	Expr        string `json:"expr"`
	Exclude     string `json:"exclude"`
}

func newFilter(id string, args string) filter {
	f := filter{
		ID: id,
	}
	f.parse(args)
	return f
}

func (f filter) String() string {
	return fmt.Sprintf("[%s] [%s %s %s]", f.ID, f.ChannelName, f.Expr, f.Exclude)
}

func (f *filter) parse(value string) {
	args := plugins.DefaultUtils.QuotationOrSpaceFields(value)

	f.ChannelName = args[0]
	if len(args) > 1 {
		f.Expr = args[1]
	}
	if len(args) > 2 {
		f.Exclude = args[2]
	}
}

func (f *filter) execute(message string) bool {
	if matched, err := regexp.MatchString(f.Expr, message); err != nil || !matched {
		return false
	}
	if len(f.Exclude) > 0 {
		if matched, err := regexp.MatchString(f.Exclude, message); err != nil || matched {
			return false
		}
	}
	return true
}
