package pkg

import (
	"context"

	"github.com/emersion/go-imap"
)

func FromAdress(address string) Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		for _, from := range msg.Envelope.From {
			if from.Address() == address {
				return true, nil
			}
		}
		return false, nil
	})
}

func FromPersonalName(name string) Rule {
	return RuleFunc(func(ctx context.Context, msg *imap.Message) (bool, error) {
		for _, from := range msg.Envelope.From {
			if from.PersonalName == name {
				return true, nil
			}
		}
		return false, nil
	})
}
