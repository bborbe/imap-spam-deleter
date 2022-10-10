package pkg

import (
	"context"
	"mime"

	"github.com/emersion/go-imap"
	charset "github.com/mantyr/go-charset/charset"
	_ "github.com/mantyr/go-charset/data"
	"github.com/pkg/errors"
)

func Subject(ctx context.Context, msg *imap.Message) (string, error) {
	if msg == nil {
		return "", errors.Errorf("msg nil")
	}
	if msg.Envelope == nil {
		return "", errors.Errorf("envelope nil")
	}
	return decodeSubject(msg.Envelope.Subject)
}

func decodeSubject(subject string) (string, error) {
	dec := new(mime.WordDecoder)
	dec.CharsetReader = charset.NewReader
	ret, err := dec.DecodeHeader(subject)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func SubjectEqual(expectedSubject string) Rule {
	return SubjectMatch(MatcherEqual(expectedSubject))

}

func SubjectPrefix(expectedPrefix string) Rule {
	return SubjectMatch(MatcherPrefix(expectedPrefix))

}

func SubjectContains(expectedSubstring string) Rule {
	return SubjectMatch(MatcherContains(expectedSubstring))
}

func SubjectMatch(matcher Matcher) Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		subject, err := Subject(ctx, msg)
		if err != nil {
			return false, err
		}
		if matcher.Match(subject) {
			return true, nil
		}
		return false, nil
	})
}
