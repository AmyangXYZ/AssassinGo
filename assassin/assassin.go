package assassin

import (
	"../logger"
)

// Assassin shadow and assassinate the target.
type Assassin struct {
	Target       string
	FuzzableURLs []string
}

// New returns a new Assassin.
func New(target string) *Assassin {
	logger.Green.Println("Target:", target)
	return &Assassin{
		Target: target,
	}
}
