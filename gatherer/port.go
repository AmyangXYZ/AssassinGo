package gatherer

import (
	"net"
	"strings"
	"time"

	"logger"
	"utils"
	"github.com/AmyangXYZ/barbarian"
	"github.com/gorilla/websocket"
)

// PortScanner scans common used ports.
// WebSocket API.
type PortScanner struct {
	mconn  *utils.MuxConn
	target string
	// tcp, syn ...
	method      string
	ports       map[string]string
	concurrency int
	timeout     int
	OpenPorts   []string
}

// NewPortScanner returns a PortScanner.
func NewPortScanner() *PortScanner {
	return &PortScanner{
		mconn:       &utils.MuxConn{},
		ports:       makePortsMap("./dict/Top100ports.txt"),
		concurrency: 100,
		timeout:     3,
	}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target, method string}
func (ps *PortScanner) Set(v ...interface{}) {
	ps.mconn.Conn = v[0].(*websocket.Conn)
	ps.target = v[1].(string)
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
	portNumbers := []string{}
	for p := range ps.ports {
		portNumbers = append(portNumbers, p)
	}
	bb := barbarian.New(ps.checkPort, ps.onResult, portNumbers, ps.concurrency)
	bb.Run()
}

func (ps *PortScanner) onResult(res interface{}) {
	port := res.(string)
	logger.Blue.Printf("%-5s -  %s \n", port, ps.ports[port])
	ret := map[string]interface{}{
		"port":    port,
		"service": ps.ports[port],
	}
	ps.mconn.Send(ret)
	ps.OpenPorts = append(ps.OpenPorts, port)
}

func (ps *PortScanner) checkPort(port string) interface{} {
	connection, err := net.DialTimeout("tcp", ps.target+":"+port, time.Duration(ps.timeout)*time.Second)
	if err == nil {
		connection.Close()
		return port
	}
	return nil
}

func makePortsMap(file string) map[string]string {
	portsMap := map[string]string{}
	for _, row := range utils.ReadFile(file) {
		x := strings.Split(row, " ")
		portsMap[x[0]] = x[1]
	}
	return portsMap
}
