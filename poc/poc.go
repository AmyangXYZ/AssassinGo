package poc

// PoC needs to implements:
// Info() -> return brief introduction
// Set() -> set params
// Run -> run the poc
// Report -> return result
type PoC interface {
	Info() string
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
