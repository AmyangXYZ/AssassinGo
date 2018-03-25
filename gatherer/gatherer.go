package gatherer

// Gatherer should impletement ...
type Gatherer interface {
	Run()
}

// Set sets the gatherers to use.
func Set(target string) []Gatherer {
	return []Gatherer{
		NewBasicInfo(target),
		NewCMSDetector(target),
		NewPortScanner(target),
	}
}
