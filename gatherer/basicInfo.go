package gatherer

import (
	"net"
	"net/http"

	"../logger"
	"github.com/gorilla/websocket"
)

// BasicInfo gathers basic information of the target.
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
// Params should be {target string}
func (bi *BasicInfo) Set(v ...interface{}) {
	bi.target = v[0].(string)
}

// Report implements Gatherer interface
func (bi *BasicInfo) Report() interface{} {
	return []string{bi.IPAddr, bi.WebServer}
}

// Run implements the Gatherer interface.
func (bi *BasicInfo) Run(conn *websocket.Conn) {
	err := bi.resolveIP()
	if err != nil {
		conn.Close()
	}
	logger.Green.Println("IP Address:", bi.IPAddr)

	err = bi.getWebServer()
	logger.Green.Println("Web Server:", bi.WebServer)

	ret := map[string]string{
		"ip":        bi.IPAddr,
		"webserver": bi.WebServer,
	}
	conn.WriteJSON(ret)
}

func (bi *BasicInfo) resolveIP() error {
	remoteAddr, err := net.ResolveIPAddr("ip", bi.target)
	if err != nil {
		logger.Red.Println(err)
		return err

	}
	bi.IPAddr = remoteAddr.String()
	return nil
}

func (bi *BasicInfo) getWebServer() error {
	resp, err := http.Head("http://" + bi.target)
	if err != nil {
		resp, err = http.Get("http://" + bi.target)
		if err != nil {
			logger.Red.Println(err)
			return err
		}
	}
	bi.WebServer = resp.Header["Server"][0]
	return nil
}
