package assassin

import (
	"strings"

	"attacker"
	"gatherer"
	"logger"
	"poc"
	"seeker"
	"utils"
)

// Daddy is used for multiple users.
type Daddy struct {
	Son     map[string]*Assassin
	Sibling map[string]*Sibling
}

// NewDaddy returns a new daddy.
func NewDaddy() *Daddy {
	return &Daddy{
		Son:     make(map[string]*Assassin),
		Sibling: make(map[string]*Sibling),
	}
}

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

// Sibling is a batch vul scanner.
type Sibling struct {
	MuxConn          utils.MuxConn
	Siblings         []*Assassin
	ExploitableHosts []string
}

// NewSiblings returns a group of Assassins for batch scan.
func NewSiblings() *Sibling {
	return &Sibling{
		Siblings: make([]*Assassin, 0),
		MuxConn:  utils.MuxConn{}}
}

// SetTargets .
func (s *Sibling) SetTargets(targets string) {
	ts := strings.Split(targets, ",")
	logger.Green.Println("Set Siblings' Targets", len(ts))
	for _, t := range ts {
		a := &Assassin{PoC: poc.Init()}
		a.Target = t
		s.Siblings = append(s.Siblings, a)
	}
}
