package gatherer

// Gatherer should impletement ...
type Gatherer interface {
	Run()
	Report() interface{}
}

// Set sets the gatherers to use.
func Set(target string) []Gatherer {
	return []Gatherer{
		NewBasicInfo(target),
		NewCMSDetector(target),
		NewPortScanner(target),
	}
}
