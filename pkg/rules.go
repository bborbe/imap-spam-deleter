package pkg

import (
	"context"

	"github.com/emersion/go-imap"
)

var rules = Or(
	SubjectContains("Bauchfett"),
	SubjectContains("Samurai-Küchenmesser"),
	SubjectContains("Erektionen"),
	SubjectContains("Erektionshilfen"),
	SubjectContains("Sexualorgane"),
	SubjectContains("Diabetes"),
	SubjectContains("erases fat"),
	SubjectContains("Kostenlose Schnelltests"),
	SubjectContains("Beischlaf"),
	SubjectContains("energiebooster"),
	SubjectContains("befreie dein wildes tier"),
	SubjectContains("Sex oder Covid"),
	SubjectContains("Sportwetten Trick zum absahnen"),
	SubjectContains("gefährlicher als Zucker"),
	And(
		SubjectContains("Scheck"),
		SubjectContains("US-Dollar"),
	),
	SubjectContains("AW: Deine Stromrechnungen gehört der Vergangenheit an"),
	SubjectContains("Bauchfett"),
	SubjectContains("kinderleicht die präzisesten dünnsten Scheiben"),
	SubjectContains("Anthelminthikum"),
	SubjectContains("Sehr geehrter Begünstigter"),
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
