package scanner

// Scanner should impletement ...
// Add your url-based scanner here.
type Scanner interface {
	Run(fuzzableURLs []string)
	Report() interface{}
}
