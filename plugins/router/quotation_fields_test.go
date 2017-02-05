package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_quotationOrSpaceFields(t *testing.T) {
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

	for i, ts := range testCases {
		as.Equal(quotationOrSpaceFields(ts.s), ts.args, "testCase [%d]", i)
	}
}
