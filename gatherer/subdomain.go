package gatherer

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/AmyangXYZ/barbarian"

	"../logger"
	"../utils"
	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/evilsocket/brutemachine"
	"github.com/gorilla/websocket"
)

// SubDomainScan brute force the dir.
// WebSocker API.
type SubDomainScan struct {
	mconn       *utils.MuxConn
	target      string
	m           *brutemachine.Machine
	wordlist    []string
	concurrency int
	wildcard    []string
	Subdomains  []string
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
		mconn:       &utils.MuxConn{},
		wordlist:    utils.ReadFile("/dict/names.txt"),
		concurrency: 20,
	}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target, concurrency int}
func (s *SubDomainScan) Set(v ...interface{}) {
	s.mconn.Conn = v[0].(*websocket.Conn)
	s.target = domainutil.Domain(v[1].(string))
}

// Report implements Gatherer interface.
func (s *SubDomainScan) Report() map[string]interface{} {
	return map[string]interface{}{
		"subdomains": s.Subdomains,
	}
}

// Run implements Gatherer interface,
func (s *SubDomainScan) Run() {
	logger.Green.Println("Enumerating Subdomain with DNS Search...")
	hasWildcard := false
	hasWildcard, s.wildcard, _ = s.detectWildcard()
	if hasWildcard {
		logger.Blue.Printf("Detected Wildcard : %v\n", s.wildcard)
	}

	bb := barbarian.New(s.DoRequest, s.OnResult, s.wordlist, s.concurrency)

	// dns lookup is easy to go out of time.
	go func() {
		time.Sleep(15 * time.Second)
		bb.Stop()
	}()

	bb.Run()
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
		if reflect.DeepEqual(addrs, s.wildcard) {
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
	s.Subdomains = append(s.Subdomains, result.hostname)
	ret := map[string]interface{}{
		"subdomain": result.hostname,
	}
	s.mconn.Send(ret)
}
