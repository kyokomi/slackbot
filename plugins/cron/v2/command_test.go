package cron

import (
	"fmt"
	"testing"
)

func TestGenerateCronID(t *testing.T) {
	fmt.Println(generateCronID())
}
