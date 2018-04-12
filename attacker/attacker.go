package attacker

import "github.com/gorilla/websocket"

// Attacker should implement ...
// Add your url-based attacker here.
type Attacker interface {
	Set(...interface{})
	Run(conn *websocket.Conn)
	Report() interface{}
}

// Init Attackers.
func Init() map[string]Attacker {
	return map[string]Attacker{
		"sqli":     NewBasicSQLi(),
		"xss":      NewXSSChecker(),
		"intruder": NewIntruder(),
	}
}
