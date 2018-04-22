// Todo: SYN scan.

package gatherer

import (
	"io/ioutil"
	"net"
	"strings"
	"time"

	"../logger"
	"../util"
	"github.com/gorilla/websocket"
)

// PortScanner scans common used ports.
// WebSocket API.
type PortScanner struct {
	mconn  *util.MuxConn
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
		mconn:           &util.MuxConn{},
		ports:           makePortsMap("/home/amyang/Projects/AssassinGo/gatherer/dict/Top100ports.txt"),
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
	portsMap := map[string]string{}
	buf, _ := ioutil.ReadFile(file)
	rows := strings.Split(string(buf), "\n")
	for _, row := range rows {
		x := strings.Split(row, " ")
		portsMap[x[0]] = x[1]
	}
	return portsMap
}
