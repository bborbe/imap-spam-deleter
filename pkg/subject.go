package pkg

import (
	"context"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/pkg/errors"
)

func Subject(ctx context.Context, msg *imap.Message) (string, error) {
	if msg == nil {
		return "", errors.Errorf("msg nil")
	}
	if msg.Envelope == nil {
		return "", errors.Errorf("envelope nil")
	}
	return msg.Envelope.Subject, nil
}

func SubjectEqual(expectedSubject string) Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		subject, err := Subject(ctx, msg)
		if err != nil {
			return false, err
		}
		if subject == expectedSubject {
			return true, nil
		}
		return false, nil
	})
}

func SubjectPrefix(prefix string) Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		subject, err := Subject(ctx, msg)
		if err != nil {
			return false, err
		}
		if strings.HasPrefix(subject, prefix) {
			return true, nil
		}
		return false, nil
	})
}

func SubjectContains(substr string) Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		subject, err := Subject(ctx, msg)
		if err != nil {
			return false, err
		}
		if strings.Contains(subject, substr) {
			return true, nil
		}
		return false, nil
	})
}
