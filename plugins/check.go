package plugins

import (
	"strings"

	"golang.org/x/net/context"
)

func CheckMessageKeyword(ctx context.Context, message string, keyword string) (bool, string) {
	if strings.Index(strings.ToLower(message), keyword) == -1 {
		return false, message
	}
	return true, message
}
