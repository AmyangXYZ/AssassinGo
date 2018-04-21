package assassin

import (
	"strings"

	"../attacker"
	"../gatherer"
	"../logger"
	"../poc"
	"../seeker"
	"../util"
)

// Assassin shadow and assassinate the target.
type Assassin struct {
	Target       string
	FuzzableURLs []string
	Gatherers    map[string]gatherer.Gatherer
	Attackers    map[string]attacker.Attacker
	Seeker       seeker.Seeker
	PoC          map[string]poc.PoC
}

// New returns a new Assassin.
func New() *Assassin {
	return &Assassin{
		Gatherers: gatherer.Init(),
		Attackers: attacker.Init(),
		PoC:       poc.Init(),
	}
}

// SetTarget .
func (a *Assassin) SetTarget(target string) {
	logger.Green.Println("Target:", target)
	a.Target = target
}

// Dad is a batch vul scanner.
type Dad struct {
	MuxConn          util.MuxConn
	Sons             []*Assassin
	ExploitableHosts []string
}

// NewDad returns Assassins' Dad.
func NewDad() *Dad {
	return &Dad{Sons: make([]*Assassin, 0)}
}

// SetTargets .
func (d *Dad) SetTargets(targets string) {
	logger.Green.Println("Set Dad's Targets")
	for _, t := range strings.Split(targets, ",") {
		son := New()
		son.Target = t
		d.Sons = append(d.Sons, son)
	}
}
