package pkg

import (
	"context"
	"time"

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
	updates := make(chan client.Update)

	c, err := client.DialTLS(a.Server, nil)
	if err != nil {
		return errors.Wrapf(err, "connected failed")
	}
	glog.V(2).Infof("Connected")

	defer func() {
		if err := c.Logout(); err != nil {
			glog.V(2).Infof("logout failed: %v", err)
		}
	}()

	if err := c.Login(a.Username, a.Password); err != nil {
		return errors.Wrapf(err, "login failed")
	}
	glog.V(2).Infof("Logged in")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := a.cleanInbox(ctx, c); err != nil {
				return errors.Wrap(err, "clean failed")
			}
			glog.V(2).Infof("cleanup completed => wait for updates")
			if err := a.waitForUpdates(ctx, c, updates); err != nil {
				return errors.Wrap(err, "update clean failed")
			}
			glog.V(2).Infof("wait for updates completed")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.NewTimer(time.Second).C:
			}
		}
	}
}

func (a *Application) cleanInbox(ctx context.Context, c *client.Client) error {
	mailboxCleaner := NewMailboxCleaner(c, a.DryRun)
	if err := mailboxCleaner.Clean(ctx, "INBOX"); err != nil {
		return err
	}
	return nil
}

func (a *Application) waitForUpdates(ctx context.Context, c *client.Client, updates chan client.Update) error {
	c.Updates = updates
	defer func() {
		c.Updates = nil
	}()

	errs := make(chan error)
	stop := make(chan struct{})
	go func() {
		errs <- c.Idle(stop, nil)
		close(errs)
	}()
	defer close(stop)

	for {

		select {
		case update := <-updates:
			switch obj := update.(type) {
			case *client.MailboxUpdate:
				glog.V(2).Infof("MailboxUpdate: %s", obj.Mailbox.Name)
				if obj.Mailbox.Name != "INBOX" {
					glog.V(2).Infof("skip mailbox %s", obj.Mailbox.Name)
					continue
				}
				return nil
			case *client.MessageUpdate:
				if obj.Message == nil {
					glog.Warning("Message is nil")
					continue
				}
				if obj.Message.Envelope == nil {
					glog.Warning("Message.Envelope is nil")
					continue
				}
				glog.V(3).Infof("MessageUpdate from %s", obj.Message.Envelope.Subject)
			case *client.StatusUpdate:
				glog.V(3).Infof("StatusUpdate %+v", obj)
			case *client.ExpungeUpdate:
				glog.V(3).Infof("ExpungeUpdate %+v", obj)
			default:
				glog.V(3).Infof("update with type %T", update)
			}
		case err := <-errs:
			glog.V(2).Infof("unexpected error: %v", err)
			return nil
		case <-ctx.Done():
			glog.V(2).Infof("ctx canceled")
			return ctx.Err()
		}
	}

}
