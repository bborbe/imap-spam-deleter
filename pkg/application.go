package pkg

import (
	"context"

	"github.com/bborbe/run"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type Application struct {
	Username string `arg:"imap-user" env:"IMAP_USER"`
	Password string `arg:"imap-password" env:"IMAP_PASSWORD" display:"length"`
	Server   string `arg:"imap-server" env:"IMAP_SERVER"`
	DryRun   bool   `arg:"dry-run" env:"DRY_RUN" default:"true"`
}

func (a *Application) Run(ctx context.Context) error {

	c, err := client.DialTLS(a.Server, nil)
	if err != nil {
		return errors.Wrapf(err, "connected failed")
	}
	glog.V(2).Infof("Connected")

	// Don't forget to logout
	defer func() {
		if err := c.Logout(); err != nil {
			glog.V(2).Infof("logout failed: %v", err)
		}
	}()

	// Login
	if err := c.Login(a.Username, a.Password); err != nil {
		return errors.Wrapf(err, "login failed")
	}
	glog.V(2).Infof("Logged in")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return errors.Wrapf(err, "select inbox failed")
	}
	glog.V(3).Infof("Flags for INBOX: %+v", mbox.Flags)

	if mbox.Messages == 0 {
		glog.V(2).Infof("empty inbox")
		return nil
	}

	from := uint32(1)
	to := mbox.Messages
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	messages := make(chan *imap.Message, 10)

	return run.CancelOnFirstError(
		ctx,
		func(ctx context.Context) error {
			//defer close(messages)
			return c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		},
		func(ctx context.Context) error {

			deleteSet := new(imap.SeqSet)

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case msg, ok := <-messages:
					if !ok {
						glog.V(2).Infof("read all message completed")

						if !a.DryRun && !deleteSet.Empty() {
							item := imap.FormatFlagsOp(imap.AddFlags, true)
							flags := []interface{}{imap.DeletedFlag}
							if err := c.Store(deleteSet, item, flags, nil); err != nil {
								return errors.Wrap(err, "store failed")
							}

							if err := c.Expunge(nil); err != nil {
								return errors.Wrap(err, "expunge failed")
							}
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

var rules = Or(
	SubjectContains("Bauchfett"),
	SubjectContains("Samurai-Küchenmesser"),
	SubjectContains("Erektionen"),
	SubjectContains("Sexualorgane"),
	SubjectContains("Diabetes"),
	SubjectContains("erases fat"),
)
