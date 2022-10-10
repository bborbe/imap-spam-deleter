package pkg

import (
	"strings"
)

type Matcher interface {
	Match(value string) bool
}

type MatcherFunc func(value string) bool

func (m MatcherFunc) Match(value string) bool {
	return m(value)
}

func MatcherContains(str string) Matcher {
	return MatcherFunc(func(value string) bool {
		return strings.Contains(value, str)
	})
}

func MatcherPrefix(str string) Matcher {
	return MatcherFunc(func(value string) bool {
		return strings.HasPrefix(value, str)
	})
}
func MatcherSuffix(str string) Matcher {
	return MatcherFunc(func(value string) bool {
		return strings.HasSuffix(value, str)
	})
}

func MatcherEqual(str string) Matcher {
	return MatcherFunc(func(value string) bool {
		return str == value
	})
}

func MatcherCleaner(matcher Matcher, cleaner StringCleaner) Matcher {
	return MatcherFunc(func(value string) bool {
		value = cleaner.Clean(value)
		return matcher.Match(value)
	})
}
