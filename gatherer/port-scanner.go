// Todo: SYN scan.

package gatherer

import (
	"io/ioutil"
	"net"
	"strings"
	"time"

	"../logger"
)

// PortScanner scans common used ports.
type PortScanner struct {
	target      string
	ports       []string
	concurrency int
	timeout     int
	OpenPorts   []string
}

// NewPortScanner returns a PortScanner.
func NewPortScanner(target string) *PortScanner {
	return &PortScanner{
		target:      target,
		ports:       readPortsFromFile("/home/amyang/Projects/AssassinGo/gatherer/Top100ports.txt"),
		concurrency: 100,
		timeout:     2,
	}
}

// Report impletements Gatherer interface
func (ps *PortScanner) Report() interface{} {
	return ps.OpenPorts
}

// Run impletements the Gatherer interface.
func (ps *PortScanner) Run() {
	logger.Green.Println("Ports Scanning...")

	blockers := make(chan bool, ps.concurrency)
	for _, port := range ps.ports {
		blockers <- true
		go ps.checkPort(port, blockers)
	}

	// Wait for all goroutines to finish.
	for i := 0; i < cap(blockers); i++ {
		blockers <- true
	}
}

func (ps *PortScanner) checkPort(port string, blocker chan bool) {
	defer func() { <-blocker }()
	connection, err := net.DialTimeout("tcp", ps.target+":"+port, time.Duration(ps.timeout)*time.Second)
	if err == nil {
		logger.Blue.Printf("%-5s -  open \n", port)
		connection.Close()
		ps.OpenPorts = append(ps.OpenPorts, port)
	}
}

func readPortsFromFile(file string) []string {
	buf, _ := ioutil.ReadFile(file)
	ports := strings.Split(string(buf), "\n")
	return ports
}
