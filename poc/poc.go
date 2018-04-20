package poc

import (
	"sync"

	"github.com/gorilla/websocket"
)

// PoC just need to implements Run().
type PoC interface {
	Set(...interface{})
	Run()
	Report() map[string]interface{}
}

// PoCMap is a poc map.
var PoCMap = map[string]PoC{
	"SeaCMSv654": NewSeaCMSv654(),
	"Drupal":     NewDrupalRCE(),
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
