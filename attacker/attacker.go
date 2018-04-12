package attacker

import "github.com/gorilla/websocket"

// Attacker should implement ...
// Add your url-based scanner here.
type Attacker interface {
	Run(fuzzableURLs []string, conn *websocket.Conn)
	Report() interface{}
}
