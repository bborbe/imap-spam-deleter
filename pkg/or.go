package pkg

import (
	"context"

	"github.com/emersion/go-imap"
)

func Or(rules ...Rule) Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		for _, rule := range rules {
			delete, err := rule.Delete(ctx, msg)
			if err != nil {
				return false, err
			}
			if delete {
				return true, nil
			}
		}
		return false, nil
	})
}
