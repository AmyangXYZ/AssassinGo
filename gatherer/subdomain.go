package gatherer

import "github.com/gorilla/websocket"

// SubDomainScan brute force the dir.
type SubDomainScan struct {
	mconn  *muxConn
	target string
}

// NewSubDomainScan returns a new SubDomainScan.
func NewSubDomainScan() *SubDomainScan {
	return &SubDomainScan{}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target, goroutinesCount int}
func (s *SubDomainScan) Set(v ...interface{}) {
	s.mconn = &muxConn{conn: v[0].(*websocket.Conn)}
	s.target = v[1].(string)
}

// Report implements Gatherer interface.
func (s *SubDomainScan) Report() interface{} {
	return ""
}

// Run implements Gatherer interface,
func (s *SubDomainScan) Run() {

}
