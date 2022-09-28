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
			glog.V(2).Infof("application shutdown")

			return ctx.Err()
		default:
			if err := a.cleanExisting(ctx, c); err != nil {
				return errors.Wrap(err, "initial clean failed")
			}
			glog.V(2).Infof("inital cleanup completed => wait for updates")
			if err := a.cleanUpdates(ctx, c); err != nil {
				return errors.Wrap(err, "update clean failed")
			}
			time.Sleep(time.Second)
		}
	}
}

func (a *Application) cleanExisting(ctx context.Context, c *client.Client) error {
	mailboxCleaner := NewMailboxCleaner(c, a.DryRun)
	if err := mailboxCleaner.Clean(ctx, "INBOX"); err != nil {
		return err
	}
	return nil
}

func (a *Application) cleanUpdates(ctx context.Context, c *client.Client) error {
	glog.V(2).Infof("cleanUpdates")
	updates := make(chan client.Update)
	c.Updates = updates
	defer func() {
		c.Updates = nil
		close(updates)
	}()

	errs := make(chan error)
	stop := make(chan struct{})
	go func() {
		errs <- c.Idle(stop, nil)
		close(errs)
	}()
	defer close(stop)

	for {
		glog.V(2).Infof("wait for update started")

		select {
		case update := <-updates:
			glog.V(2).Infof("update with type %T", update)
			switch obj := update.(type) {
			case *client.MessageUpdate:
				glog.V(2).Infof("MessageUpdate from %s", obj.Message.Envelope.Subject)

			case *client.MailboxUpdate:
				glog.V(2).Infof("MailboxUpdate: %s", obj.Mailbox.Name)
				if obj.Mailbox.Name != "INBOX" {
					glog.V(2).Infof("skip mailbox %s", obj.Mailbox.Name)
					continue
				}
				return nil
			case *client.StatusUpdate:
				glog.V(2).Infof("StatusUpdate %+v", obj)
			case *client.ExpungeUpdate:
				glog.V(2).Infof("ExpungeUpdate %+v", obj)
			}
			glog.V(2).Infof("handle update completed")
		case err := <-errs:
			glog.V(2).Infof("unexpected error: %v", err)
		case <-ctx.Done():
			glog.V(2).Infof("ctx canceled")
			return ctx.Err()
		}
	}

}
