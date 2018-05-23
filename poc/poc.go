package poc

// Intro is a brief intro of the PoC.
type Intro struct {
	// CVE ID
	ID string `json:"id"`
	// RCE or SQLi or ...
	Type string `json:"type"`
	// Breif Intro
	Text string `json:"text"`
	// php or python or ...
	Platform string `json:"platform"`
	// 2016-3-11
	Date string `json:"date"`
	// http://....
	Reference string `json:"reference"`
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
