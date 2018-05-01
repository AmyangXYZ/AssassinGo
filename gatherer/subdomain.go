// Adapted from https://github.com/evilsocket/dnssearch

package gatherer

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"reflect"

	"../logger"
	"../util"
	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/evilsocket/brutemachine"
	"github.com/gorilla/websocket"
)

// SubDomainScan brute force the dir.
type SubDomainScan struct {
	mconn      *util.MuxConn
	target     string
	subdomains []string
}

// Result to show what we've found
type Result struct {
	hostname string
	addrs    []string
	txts     []string
	cname    string // Per RFC, there should only be one CNAME
}

// NewSubDomainScan returns a new SubDomainScan.
func NewSubDomainScan() *SubDomainScan {
	return &SubDomainScan{
		mconn: &util.MuxConn{},
	}
}

var (
	m *brutemachine.Machine

	wordlist  = "/dict/names.txt" // Wordlist file to use for enumeration.
	consumers = 8                 // Number of concurrent consumers.
	forceTld  = true              // Extract top level from provided domain

	wildcard []string
)

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target, goroutinesCount int}
func (s *SubDomainScan) Set(v ...interface{}) {
	s.mconn.Conn = v[0].(*websocket.Conn)
	s.target = domainutil.Domain(v[1].(string))
	hasWildcard := false

	hasWildcard, wildcard, _ = s.detectWildcard()

	if hasWildcard {
		logger.Blue.Printf("Detected Wildcard : %v\n\n", wildcard)
	}
}

// Report implements Gatherer interface.
func (s *SubDomainScan) Report() map[string]interface{} {
	return map[string]interface{}{
		"subdomains": s.subdomains,
	}
}

// Run implements Gatherer interface,
func (s *SubDomainScan) Run() {
	logger.Green.Println("Enumerating Subdomain with DNS Search...")

	m = brutemachine.New(consumers, wordlist, s.DoRequest, s.OnResult)
	if err := m.Start(); err != nil {
		logger.Red.Println(err)
		return
	}

	m.Wait()
}

// Lookup a random host to determine if a wildcard A record exists
// Adapted from https://github.com/jrozner/sonar/blob/master/wildcard.go
func (s *SubDomainScan) detectWildcard() (bool, []string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return false, nil, err
	}

	domain := fmt.Sprintf("%s.%s", hex.EncodeToString(bytes), s.target)

	answers, err := net.LookupHost(domain)
	if err != nil {
		if asserted, ok := err.(*net.DNSError); ok && asserted.Err == "no such host" {
			return false, nil, nil
		}
		return false, nil, err
	}

	return true, answers, nil
}

// DoRequest actually handles the DNS lookups
func (s *SubDomainScan) DoRequest(sub string) interface{} {
	hostname := fmt.Sprintf("%s.%s", sub, s.target)
	thisresult := Result{}

	if addrs, err := net.LookupHost(hostname); err == nil {
		if reflect.DeepEqual(addrs, wildcard) {
			// This is likely a wildcard entry, skip it
			return nil
		}
		thisresult.hostname = hostname
		thisresult.addrs = addrs
	}

	if thisresult.hostname == "" {
		return nil
	}

	return thisresult
}

// OnResult prints out the results of a lookup
func (s *SubDomainScan) OnResult(res interface{}) {
	result, ok := res.(Result)
	if !ok {
		logger.Red.Printf("Error while converting result.\n")
		return
	}

	logger.Blue.Println(result.hostname)
	s.subdomains = append(s.subdomains, result.hostname)
	ret := map[string]interface{}{
		"subdomain": result.hostname,
	}
	s.mconn.Send(ret)
}
