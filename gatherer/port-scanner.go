// Todo: SYN scan.

package gatherer

import (
	"net"
	"strings"
	"time"

	"../logger"
	"../utils"
	"github.com/gorilla/websocket"
)

// PortScanner scans common used ports.
// WebSocket API.
type PortScanner struct {
	mconn  *utils.MuxConn
	target string
	// tcp, syn ...
	method          string
	ports           map[string]string
	goroutinesCount int
	timeout         int
	OpenPorts       []string
}

// NewPortScanner returns a PortScanner.
func NewPortScanner() *PortScanner {
	return &PortScanner{
		mconn:           &utils.MuxConn{},
		ports:           makePortsMap("./dict/Top100ports.txt"),
		goroutinesCount: 100,
		timeout:         3,
	}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target, method string}
func (ps *PortScanner) Set(v ...interface{}) {
	ps.mconn.Conn = v[0].(*websocket.Conn)
	ps.target = v[1].(string)
	ps.method = v[2].(string)
}

// Report implements Gatherer interface
func (ps *PortScanner) Report() map[string]interface{} {
	return map[string]interface{}{
		"ports": ps.OpenPorts,
	}
}

// Run implements the Gatherer interface.
func (ps *PortScanner) Run() {
	logger.Green.Println("Ports Scanning...")
	ps.OpenPorts = []string{}

	blockers := make(chan struct{}, ps.goroutinesCount)
	for port := range ps.ports {
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
		logger.Blue.Printf("%-5s -  %s \n", port, ps.ports[port])
		connection.Close()
		ret := map[string]interface{}{
			"port":    port,
			"service": ps.ports[port],
		}
		ps.mconn.Send(ret)
		ps.OpenPorts = append(ps.OpenPorts, port)
	}
}

func makePortsMap(file string) map[string]string {
	data, err := utils.ReadFile(file)
	if err != nil {
		logger.Red.Println(err)
		return map[string]string{}
	}

	portsMap := map[string]string{}
	for _, row := range data {
		x := strings.Split(row, " ")
		portsMap[x[0]] = x[1]
	}
	return portsMap
}
