package poc

// Intro is a brief intro of the PoC.
type Intro struct {
	// CVE ID
	ID string
	// RCE or SQLi or ...
	Type string
	// Breif Intro
	Text string
	// php or python or ...
	Platform string
	// 2016-3-11
	Date string
	// http://....
	Reference string
}

// PoC needs to implements:
// Info() -> return brief introduction
// Set() -> set params
// Run -> run the poc
// Report -> return result
type PoC interface {
	Info() Intro
	Set(...interface{})
	Run()
	Report() map[string]interface{}
}

// Init PoC
func Init() map[string]PoC {
	return map[string]PoC{
		"seacms-v654-rce": NewSeaCMSv654(),
		"drupal-rce":      NewDrupalRCE(),
	}
}
