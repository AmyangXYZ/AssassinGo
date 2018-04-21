package attacker

import (
	"io/ioutil"
	"net/http"
	"strings"

	"../logger"
	"../util"
	"github.com/gorilla/websocket"
)

// XSSChecker checks XSS vuls.
// WebSocket API.
type XSSChecker struct {
	mconn         *util.MuxConn
	fuzzableURLs  []string
	payload       string
	InjectableURL []string
}

// NewXSSChecker returns a XSS Checker.
func NewXSSChecker() *XSSChecker {
	return &XSSChecker{payload: `<svg/onload=alert(1)>`}
}

// Set implements Attacker interface.
// Params should be {conn *websocket.Conn, fuzzableURLs []string}
func (x *XSSChecker) Set(v ...interface{}) {
	x.mconn.Conn = v[0].(*websocket.Conn)
	x.fuzzableURLs = v[1].([]string)
}

// Report implements Attacker interface.
func (x *XSSChecker) Report() map[string]interface{} {
	return map[string]interface{}{
		"xss_urls": x.InjectableURL,
	}
}

// Run implements Attacker interface.
func (x *XSSChecker) Run() {
	logger.Green.Println("Basic XSS Checking...")
	x.InjectableURL = []string{}

	blockers := make(chan bool, len(x.fuzzableURLs))
	for _, URL := range x.fuzzableURLs {
		blockers <- true
		go x.check(URL, blockers)
	}

	// Wait for all goroutines to finish.
	for i := 0; i < cap(blockers); i++ {
		blockers <- true
	}

	if len(x.InjectableURL) == 0 {
		logger.Blue.Println("no xss vuls found")
	}
}

func (x *XSSChecker) check(URL string, blocker chan bool) {
	defer func() { <-blocker }()
	body := x.fetch(URL + x.payload)
	if strings.Contains(body, x.payload) {
		logger.Blue.Println(URL + x.payload)
		ret := map[string]string{
			"xss_url": URL,
		}
		x.mconn.Send(ret)
		x.InjectableURL = append(x.InjectableURL, URL)
	}
}

func (x *XSSChecker) fetch(URL string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body)
}
