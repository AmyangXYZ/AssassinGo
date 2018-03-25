package gatherer

import (
	"net"
	"net/http"

	"../logger"
)

// BasicInfo gathers basic information of the target.
type BasicInfo struct {
	target    string
	IPAddr    string
	WebServer string
}

// NewBasicInfo returns a basicInfo gatherer.
func NewBasicInfo(target string) *BasicInfo {
	return &BasicInfo{target: target}
}

// Run impletements the Gatherer interface.
func (bi *BasicInfo) Run() {
	bi.IPAddr = bi.resolveIP()
	logger.Green.Println("IP Address:", bi.IPAddr)

	bi.WebServer = bi.getWebServer()
	logger.Green.Println("Web Server:", bi.WebServer)
}

func (bi *BasicInfo) resolveIP() string {
	remoteAddr, err := net.ResolveIPAddr("ip", bi.target)
	if err != nil {
		logger.Red.Fatal(err)
	}
	return remoteAddr.String()
}

func (bi *BasicInfo) getWebServer() string {
	resp, err := http.Head("http://" + bi.target)
	if err != nil {
		resp, err = http.Get("http://" + bi.target)
		if err != nil {
			logger.Red.Fatal(err)
		}
	}
	return resp.Header["Server"][0]
}
