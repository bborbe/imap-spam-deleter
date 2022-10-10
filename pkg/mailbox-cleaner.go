package pkg

import (
	"context"

	"github.com/bborbe/run"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type MailboxCleaner interface {
	Clean(ctx context.Context, name string) error
}

func NewMailboxCleaner(client *client.Client, dryRun bool) MailboxCleaner {
	return &mailboxCleaner{
		client: client,
		dryRun: dryRun,
	}
}

type mailboxCleaner struct {
	client *client.Client
	dryRun bool
}

func (m *mailboxCleaner) Clean(ctx context.Context, name string) error {
	glog.V(2).Infof("clean started")
	defer glog.V(2).Infof("clean completed")
	mbox, err := m.client.Select("INBOX", false)
	if err != nil {
		return errors.Wrapf(err, "select inbox failed")
	}
	glog.V(3).Infof("Flags for INBOX: %+v", mbox.Flags)
	if mbox.Messages == 0 {
		glog.V(2).Infof("inbox is empty => skip")
		return nil
	}
	glog.V(2).Infof("handle %d messages", mbox.Messages)

	from := uint32(1)
	to := mbox.Messages
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	messages := make(chan *imap.Message, 10)

	return run.CancelOnFirstError(
		ctx,
		func(ctx context.Context) error {
			glog.V(2).Infof("fetch started")
			defer glog.V(2).Infof("fetch completed")
			//defer close(messages)
			return m.client.Fetch(
				seqset,
				[]imap.FetchItem{imap.FetchEnvelope},
				messages,
			)
		},
		func(ctx context.Context) error {

			deleteSet := new(imap.SeqSet)
			rules := BuildRules()
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case msg, ok := <-messages:
					if !ok {
						glog.V(2).Infof("check all message completed")

						if !m.dryRun && !deleteSet.Empty() {
							item := imap.FormatFlagsOp(imap.AddFlags, true)
							flags := []interface{}{imap.DeletedFlag}
							if err := m.client.Store(deleteSet, item, flags, nil); err != nil {
								return errors.Wrap(err, "store failed")
							}

							if err := m.client.Expunge(nil); err != nil {
								return errors.Wrap(err, "expunge failed")
							}
							glog.V(2).Infof("delete spam completed")
						}

						return nil
					}
					delete, err := rules.Delete(ctx, msg)
					if err != nil {
						return errors.Wrapf(err, "rules failed")
					}

					if delete {
						glog.V(2).Infof("DELETE msg %d subject: '%s'", msg.SeqNum, msg.Envelope.Subject)
						deleteSet.AddNum(msg.SeqNum)
					} else {
						glog.V(2).Infof("KEEP msg %d subject: '%s'", msg.SeqNum, msg.Envelope.Subject)
					}
				}
			}
		},
	)
}
