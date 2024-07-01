package gatherer

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"logger"
)

// BasicInfo gathers basic information of the target.
// AJAX API.
type BasicInfo struct {
	target                  string
	IPAddr                  string
	WebServer               string
	ClickJackingProtection  bool
	ContentSecurityPolicy   bool
	XContentTypeOptions     bool
	StrictTransportSecurity bool
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
func (bi *BasicInfo) Report() map[string]interface{} {
	return map[string]interface{}{
		"ip":                        bi.IPAddr,
		"webserver":                 bi.WebServer,
		"click_jacking_protection":  bi.ClickJackingProtection,
		"content_security_policy":   bi.ContentSecurityPolicy,
		"x_content_type_options":    bi.XContentTypeOptions,
		"strict_transport_security": bi.StrictTransportSecurity,
	}
}

// Run implements the Gatherer interface.
func (bi *BasicInfo) Run() {
	bi.resolveIP()
	logger.Green.Println("IP Address:", bi.IPAddr)

	bi.retrieveHeader()
	logger.Green.Println("Web Server:", bi.WebServer)
	logger.Green.Println("Click Jacking Protection:", bi.ClickJackingProtection)
	logger.Green.Println("Content Security Policy:", bi.ContentSecurityPolicy)
	logger.Green.Println("X Content Type Options:", bi.XContentTypeOptions)
	logger.Green.Println("Strict Transport Security:", bi.StrictTransportSecurity)
}

func (bi *BasicInfo) resolveIP() {
	t := bi.target
	if strings.Contains(bi.target, ":") {
		t = strings.Split(bi.target, ":")[0]
	}
	remoteAddr, err := net.ResolveIPAddr("ip", t)
	if err != nil {
		logger.Red.Println(err)
		return
	}
	bi.IPAddr = remoteAddr.String()
}

func (bi *BasicInfo) retrieveHeader() {
	resp, err := http.Get("http://" + bi.target)
	if err != nil {
		logger.Red.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if len(resp.Header["X-Frame-Options"]) > 0 {
		bi.ClickJackingProtection = true
	}
	if len(resp.Header["Content-Security-Policy"]) > 0 ||
		strings.Contains(string(body), `http-equiv="Content-Security-Policy"`) {
		bi.ContentSecurityPolicy = true
	}
	if len(resp.Header["X-Content-Type-Options"]) > 0 {
		bi.XContentTypeOptions = true
	}
	if len(resp.Header["Strict-Transport-Secruity"]) > 0 {
		bi.StrictTransportSecurity = true
	}
	bi.WebServer = resp.Header["Server"][0]
}
