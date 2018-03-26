package scanner

// Scanner should impletement ...
// Add your url-based scanner here.
type Scanner interface {
	Run(fuzzableURLs []string)
}

// Set sets the scanners to use.
func Set() []Scanner {
	return []Scanner{
		NewBasicSQli(),
		NewXSSChecker(),
	}
}
