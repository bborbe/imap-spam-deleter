package pkg

import (
	"context"

	"github.com/emersion/go-imap"
)

type Rule interface {
	Delete(ctx context.Context, msg *imap.Message) (delete bool, err error)
}

type RuleFunc func(ctx context.Context, msg *imap.Message) (bool, error)

func (r RuleFunc) Delete(ctx context.Context, msg *imap.Message) (bool, error) {
	return r(ctx, msg)
}

func True() Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		return true, nil
	})
}

func False() Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		return false, nil
	})
}
