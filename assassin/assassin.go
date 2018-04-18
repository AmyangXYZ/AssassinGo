package assassin

import (
	"../attacker"
	"../gatherer"
	"../logger"
	"../poc"
	"../seeker"
)

// Assassin shadow and assassinate the target.
type Assassin struct {
	Target       string
	FuzzableURLs []string
	Gatherers    map[string]gatherer.Gatherer
	Attackers    map[string]attacker.Attacker
	Seeker       seeker.Seeker
	seekerResult []string
	PoC          poc.PoC
}

// New returns a new Assassin.
func New(target string) *Assassin {
	logger.Green.Println("The Assassin is Coming...")
	logger.Green.Println("Target:", target)
	return &Assassin{
		Target:    target,
		Gatherers: gatherer.Init(),
		Attackers: attacker.Init(),
	}
}
