package utils

import (
	"bufio"
	"os"
	"sync"

	"../logger"
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

// ReadFile reads file in lines.
func ReadFile(f string) (data []string) {
	b, err := os.Open(f)
	if err != nil {
		logger.Red.Println(err)
		return
	}
	defer b.Close()
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return
}
