package router

import (
	"strings"
	"unicode"
)

var quotationsReplacer = strings.NewReplacer("'", "", `"`, "")

func quotationOrSpaceFields(s string) []string {
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
		args[i] = quotationsReplacer.Replace(args[i])
	}
	return args
}
