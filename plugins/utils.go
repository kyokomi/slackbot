package plugins

import "strings"

func CheckMessageKeyword(message string, keyword string) (bool, string) {
	if strings.Index(strings.ToLower(message), keyword) == -1 {
		return false, message
	}
	return true, message
}
