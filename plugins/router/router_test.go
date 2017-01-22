package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin_filter(t *testing.T) {
	as := assert.New(t)
	type testCase struct {
		command     string
		message     string
		channelName string
		isMatched   bool
	}
	testCases := []testCase{
		{
			command:     `#random 'Create Task.*'`,
			message:     "Create Task \n aaaaaaaaaaaaaaaaaa",
			channelName: "#random",
			isMatched:   true,
		},
		{
			command:     `#random`,
			message:     "Create Task \n aaaaaaaaaaaaaaaaaa",
			channelName: "#random",
			isMatched:   true,
		},
		{
			command:     `#random hogehoge`,
			message:     "Create Task \n aaaaaaaaaaaaaaaaaa",
			channelName: "",
			isMatched:   false,
		},
		{
			command:     `#random "Create Task.*" ".*SubTask.*`,
			message:     "Create Task \n aaaaaaaaaaaaaaaaaa\nSubTask: hogehoge",
			channelName: "",
			isMatched:   false,
		},
	}

	p := &plugin{}
	for i, ts := range testCases {
		channelName, isMatched := p.filter(ts.command, ts.message)
		as.Equal(channelName, ts.channelName, "testCase [%d]", i)
		as.Equal(isMatched, ts.isMatched, "testCase [%d]", i)
	}
}
