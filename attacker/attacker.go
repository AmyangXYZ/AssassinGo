package attacker

// Attacker should implement ...
// Add your url-based attacker here.
type Attacker interface {
	Set(...interface{})
	Run()
	Report() map[string]interface{}
}

// Init Attackers.
func Init() map[string]Attacker {
	return map[string]Attacker{
		"crawler":  NewCrawler(),
		"sqli":     NewBasicSQLi(),
		"xss":      NewXSSChecker(),
		"intruder": NewIntruder(),
	}
}
