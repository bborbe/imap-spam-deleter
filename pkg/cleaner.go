package pkg

import (
	"regexp"
	"strings"
)

var removeSpecialChars = regexp.MustCompile(`[^a-zA-Z0-9]+`)

type StringCleaner interface {
	Clean(str string) string
}

type StringCleanerList []StringCleaner

func (s StringCleanerList) Clean(str string) string {
	for _, ss := range s {
		str = ss.Clean(str)
	}
	return str
}

type StringCleanerFunc func(str string) string

func (s StringCleanerFunc) Clean(str string) string {
	return s(str)
}

func RemoveSpecialChars() StringCleaner {
	return StringCleanerFunc(func(str string) string {
		return removeSpecialChars.ReplaceAllString(str, "")
	})
}

func ToLower() StringCleaner {
	return StringCleanerFunc(func(str string) string {
		return strings.ToLower(str)
	})
}
