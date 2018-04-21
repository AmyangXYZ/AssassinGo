package util

import (
	"sync"

	"github.com/gorilla/websocket"
)

// MuxConn is a websocket.Conn with Mutex.
type MuxConn struct {
	Conn *websocket.Conn
	mu   sync.Mutex
}

// Send sends JSON.
func (m *MuxConn) Send(v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Conn.WriteJSON(v)
}

// Signal pass signal from client to server.
type Signal struct {
	Stop int
}
