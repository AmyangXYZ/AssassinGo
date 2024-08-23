package gatherer

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"assassingo/logger"
)

// CFBypass finds the real IP behind cloudflare.
// AJAX API
type CFBypass struct {
	target string
	re     *regexp.Regexp
	RealIP string
}

// NewCFBypass returns a new CFBypass
func NewCFBypass() *CFBypass {
	return &CFBypass{re: regexp.MustCompile(`(\d*\.\d*\.\d*\.\d*)`)}
}

// Set implements Gatherer interface.
// Params should be {target string}
func (cf *CFBypass) Set(v ...interface{}) {
	cf.target = v[0].(string)
}

// Report implements Gatherer interface
func (cf *CFBypass) Report() map[string]interface{} {
	return map[string]interface{}{
		"real_ip": cf.RealIP}
}

// Run implements Gatherer interaface.
func (cf *CFBypass) Run() {
	resp, err := http.Post("http://www.crimeflare.us:82/cgi-bin/cfsearch.cgi",
		"application/x-www-form-urlencoded",
		strings.NewReader("cfS="+cf.target))
	if err != nil {
		logger.Red.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	realIP := cf.re.FindAllString(string(body), 1)
	if len(realIP) > 0 {
		cf.RealIP = realIP[0]
		logger.Green.Println("RealIP:", cf.RealIP)
		return
	}
}
