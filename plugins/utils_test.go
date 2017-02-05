package plugins_test

import (
	"testing"

	"github.com/kyokomi/slackbot/plugins"
	"github.com/stretchr/testify/assert"
)

func TestUtils_quotationOrSpaceFields(t *testing.T) {
	as := assert.New(t)
	type testCase struct {
		s    string
		args []string
	}
	testCases := []testCase{
		{
			s:    `#random 'Create Task.*'`,
			args: []string{"#random", "Create Task.*"},
		},
		{
			s:    `#random "Create Task.*"`,
			args: []string{"#random", "Create Task.*"},
		},
	}

	ut := plugins.NewUtils(nil)
	for i, ts := range testCases {
		as.Equal(ut.QuotationOrSpaceFields(ts.s), ts.args, "testCase [%d]", i)
	}
}
