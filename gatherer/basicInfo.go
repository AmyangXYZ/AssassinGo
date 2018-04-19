package gatherer

import (
	"net"
	"net/http"

	"../logger"
)

// BasicInfo gathers basic information of the target.
// AJAX API.
type BasicInfo struct {
	target    string
	IPAddr    string
	WebServer string
}

// NewBasicInfo returns a basicInfo gatherer.
func NewBasicInfo() *BasicInfo {
	return &BasicInfo{}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target string}
func (bi *BasicInfo) Set(v ...interface{}) {
	bi.target = v[0].(string)
}

// Report implements Gatherer interface
func (bi *BasicInfo) Report() map[string]interface{} {
	return map[string]interface{}{
		"ip":        bi.IPAddr,
		"webserver": bi.WebServer,
	}
}

// Run implements the Gatherer interface.
func (bi *BasicInfo) Run() {
	bi.resolveIP()
	logger.Green.Println("IP Address:", bi.IPAddr)

	bi.getWebServer()
	logger.Green.Println("Web Server:", bi.WebServer)
}

func (bi *BasicInfo) resolveIP() {
	remoteAddr, err := net.ResolveIPAddr("ip", bi.target)
	if err != nil {
		logger.Red.Println(err)
	}
	bi.IPAddr = remoteAddr.String()
}

func (bi *BasicInfo) getWebServer() {
	resp, err := http.Head("http://" + bi.target)
	if err != nil {
		resp, err = http.Get("http://" + bi.target)
		if err != nil {
			logger.Red.Println(err)
		}
	}
	bi.WebServer = resp.Header["Server"][0]
}
