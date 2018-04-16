package gatherer

import (
	"io/ioutil"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// Gatherer should implement ...
type Gatherer interface {
	Set(...interface{})
	Run()
	Report() interface{}
}

// Init Gatherers.
func Init() map[string]Gatherer {
	return map[string]Gatherer{
		"basicInfo": NewBasicInfo(),
		"whois":     NewWhois(),
		"cms":       NewCMSDetector(),
		"port":      NewPortScanner(),
		"tracert":   NewTracer(),
		"dirb":      NewDirBruter(),
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

func readPayloadsFromFile(file string) []string {
	buf, _ := ioutil.ReadFile(file)
	p := strings.Split(string(buf), "\n")
	return p
}
