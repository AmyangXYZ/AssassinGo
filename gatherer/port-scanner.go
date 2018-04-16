// Todo: SYN scan.

package gatherer

import (
	"io/ioutil"
	"net"
	"strings"
	"time"

	"../logger"
	"github.com/gorilla/websocket"
)

// PortScanner scans common used ports.
type PortScanner struct {
	mconn  *muxConn
	target string
	// tcp, syn ...
	method          string
	ports           []string
	goroutinesCount int
	timeout         int
	OpenPorts       []string
}

// NewPortScanner returns a PortScanner.
func NewPortScanner() *PortScanner {
	return &PortScanner{
		ports:           readPortsFromFile("./gatherer/Top100ports.txt"),
		goroutinesCount: 100,
		timeout:         2,
	}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target, method string}
func (ps *PortScanner) Set(v ...interface{}) {
	ps.mconn = &muxConn{conn: v[0].(*websocket.Conn)}
	ps.target = v[1].(string)
	ps.method = v[2].(string)
}

// Report implements Gatherer interface
func (ps *PortScanner) Report() interface{} {
	return ps.OpenPorts
}

// Run implements the Gatherer interface.
func (ps *PortScanner) Run() {
	logger.Green.Println("Ports Scanning...")

	blockers := make(chan struct{}, ps.goroutinesCount)
	for _, port := range ps.ports {
		blockers <- struct{}{}
		go ps.checkPort(port, blockers)
	}

	// Wait for all goroutines to finish.
	for i := 0; i < cap(blockers); i++ {
		blockers <- struct{}{}
	}
}

func (ps *PortScanner) checkPort(port string, blocker chan struct{}) {
	defer func() { <-blocker }()
	connection, err := net.DialTimeout("tcp", ps.target+":"+port, time.Duration(ps.timeout)*time.Second)
	if err == nil {
		logger.Blue.Printf("%-5s -  open \n", port)
		connection.Close()
		ret := map[string]interface{}{
			"port":    port,
			"service": "http",
		}
		ps.mconn.send(ret)
		ps.OpenPorts = append(ps.OpenPorts, port)
	}
}

func readPortsFromFile(file string) []string {
	buf, _ := ioutil.ReadFile(file)
	ports := strings.Split(string(buf), "\n")
	return ports
}
