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
	SubjectContains("Cannabisöl"),
	SubjectContains("Potenzschwaeche"),
	SubjectContains("Haltungskorrektur der Wirbelsaule"),
	SubjectContains("Kostenlose Schnelltests"),
	SubjectContains("kostbaren flaschen"),
	SubjectContains("Covid impotent"),
	SubjectContains("Potenzpillen"),
	SubjectContains("Papillome"),
	SubjectContains("Hämorrhoiden"),
	FromPersonalName("Empfohlen von Pro7"),
	SubjectContains("Generikas rezeptfrei"),
	SubjectContains("Dear Beneficiary"),
	SubjectContains("Banktransfer"),
	SubjectContains("US-Dollar Erbschaft"),
	And(
		SubjectContains("Sparkasse"),
		SubjectContains("Zwei-Faktor-Authentifizierung"),
	),
	SubjectContains("Der Sportwagen unter den SUV"),
	SubjectContains("Testosteronspiegel"),
	SubjectContains("Du zahlst entscheidend zu viel für Strom"),
	SubjectContains("Fettverbrenner"),
	SubjectContains("g�nstigen Bank-�berweisung"),
	SubjectContains("Wattestäbchen"),
	SubjectContains("P-otenzmittel"),
	SubjectContains("rezeptf:rei"),
	SubjectContains("blutjunges Ding verfuehren"),
	SubjectContains("ohne Rezept via Versand"),
	SubjectContains("relief at any age"),
	SubjectContains("Sicherer Gewichtsverlust"),
	SubjectContains("Add this to your water ASAP"),
	SubjectContains("Attn: Beneficiary"),
	SubjectContains("women’s hair loss"),
	SubjectContains("Good day Dear"),
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
