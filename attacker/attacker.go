package attacker

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Attacker should implement ...
// Add your url-based attacker here.
type Attacker interface {
	Set(...interface{})
	Run()
	Report() interface{}
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

type muxConn struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (m *muxConn) send(v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.conn.WriteJSON(v)
}

type signal struct {
	Stop int
}
