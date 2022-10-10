package pkg

import (
	"context"

	"github.com/emersion/go-imap"
)

func BuildRules() Rule {

	rules := Rules{
		And(
			SubjectContains("Scheck"),
			SubjectContains("US-Dollar"),
		),
		And(
			SubjectContains("Sparkasse"),
			SubjectContains("Zwei-Faktor-Authentifizierung"),
		),
		FromPersonalName("Empfohlen von Pro7"),
	}
	var subjectContains = []string{
		"Add this to your water ASAP",
		"Anthelminthikum",
		"Attn: Beneficiary",
		"Deine Stromrechnungen gehört der Vergangenheit an",
		"Banktransfer",
		"Bauchfett",
		"Bauchfett",
		"befreie dein wildes tier",
		"Beischlaf",
		"blutjunges Ding verfuehren",
		"Cannabisöl",
		"Covid impotent",
		"Dear Beneficiary",
		"Der Sportwagen unter den SUV",
		"Diabetes",
		"Du zahlst entscheidend zu viel für Strom",
		"energiebooster",
		"erases fat",
		"Erektionen",
		"Erektionshilfen",
		"Fettverbrenner",
		"gefährlicher als Zucker",
		"Generikas rezeptfrei",
		"Good day Dear",
		"g�nstigen Bank-�berweisung",
		"Haltungskorrektur der Wirbelsaule",
		"Hämorrhoiden",
		"kinderleicht die präzisesten dünnsten Scheiben",
		"kostbaren flaschen",
		"Kostenlose Schnelltests",
		"ohne Rezept via Versand",
		"P-otenzmittel",
		"Papillome",
		"Potenzpillen",
		"Potenzschwaeche",
		"relief at any age",
		"rezeptfrei",
		"Samurai-Küchenmesser",
		"Sehr geehrter Begünstigter",
		"Sex oder Covid",
		"Sexualorgane",
		"Sicherer Gewichtsverlust",
		"Sportwetten Trick zum absahnen",
		"Testosteronspiegel",
		"Testpflicht wieder einführen",
		"US-Dollar Erbschaft",
		"Wattestäbchen",
		"women’s hair loss",
		"Sonntags-Gruß von mir",
	}
	cleaner := StringCleanerList{
		ToLower(),
		RemoveSpecialChars(),
	}
	for _, subjectContain := range subjectContains {
		rules = append(rules, SubjectMatch(
			MatcherCleaner(MatcherContains(cleaner.Clean(subjectContain)), cleaner),
		))
	}

	return Or(rules...)
}

type Rules []Rule
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
