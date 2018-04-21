package poc

// PoC just need to implements Run().
type PoC interface {
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
