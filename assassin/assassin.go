package assassin

import (
	"../crawler"
	"../gatherer"
	"../logger"
	"../poc"
	"../scanner"
)

// Assassin shadow and assassinate the target.
type Assassin struct {
	Target       string
	Crawler      crawler.Crawler
	Gatherers    map[string]gatherer.Gatherer
	FuzzableURLs []string
	Scanners     []scanner.Scanner
	POC          poc.POC
}

// New returns a new Assassin.
func New(target string) *Assassin {
	logger.Green.Println("Target:", target)
	return &Assassin{
		Target:    target,
		Crawler:   crawler.NewCrawler(target, 4),
		Gatherers: make(map[string]gatherer.Gatherer),
		Scanners:  scanner.Set(),
	}
}

// Shadow shadows the target and gathers the information.
// Run Gatherers here.
func (a *Assassin) Shadow() {
	for _, gatherer := range a.Gatherers {
		gatherer.Run()
	}

	var emails []string
	emails, a.FuzzableURLs = a.Crawler.Run()
	if len(emails) == 0 {
		logger.Green.Println("No Related E-Mails Found.")
	}
	logger.Green.Println("Related E-Mails Found:")
	for _, m := range emails {
		logger.Blue.Println(m)
	}
}

// Attack attacks the target.
// Run Scanners here.
func (a *Assassin) Attack() {
	for _, scanner := range a.Scanners {
		scanner.Run(a.FuzzableURLs)
	}
}

// Assassinate kills the target.
// Run your own POC here.
func (a *Assassin) Assassinate(POC poc.POC) {
	a.POC = POC
	a.POC.Run(a.Target)
}
