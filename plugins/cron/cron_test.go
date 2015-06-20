package cron

import (
	"fmt"
	"testing"
)

func TestCronCommand(t *testing.T) {
	command := `cron add */1 * * * * * hogehoge`

	c := CronCommand{}
	if err := c.Scan(command); err != nil {
		t.Errorf("error %s", err)
	}

	if command != fmt.Sprintf("%s", c) {
		t.Errorf("error \n%s\n%s", command, c.String())
	}
}
