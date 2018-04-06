package scanner

import (
	"io/ioutil"
	"net/http"
	"strings"

	"../logger"
	"github.com/gorilla/websocket"
)

// XSSChecker checks XSS vuls.
type XSSChecker struct {
	payload       string
	InjectableURL []string
}

// NewXSSChecker returns a XSS Checker.
func NewXSSChecker() *XSSChecker {
	return &XSSChecker{payload: `<svg/onload=alert(1)>`}
}

// Report impletements Scanner interface.
func (x *XSSChecker) Report() interface{} {
	return x.InjectableURL
}

// Run impletements Scanner interface.
func (x *XSSChecker) Run(fuzzableURLs []string, conn *websocket.Conn) {
	logger.Green.Println("Basic XSS Checking...")

	blockers := make(chan bool, len(fuzzableURLs))
	for _, URL := range fuzzableURLs {
		blockers <- true
		go x.check(URL, blockers, conn)
	}

	// Wait for all goroutines to finish.
	for i := 0; i < cap(blockers); i++ {
		blockers <- true
	}

	if len(x.InjectableURL) == 0 {
		logger.Blue.Println("no xss vuls found")
	}
}

func (x *XSSChecker) check(URL string, blocker chan bool, conn *websocket.Conn) {
	defer func() { <-blocker }()
	body := x.fetch(URL + x.payload)
	if strings.Contains(body, x.payload) {
		logger.Blue.Println(URL + x.payload)
		conn.WriteJSON(URL)
		x.InjectableURL = append(x.InjectableURL, URL)
	}
}

func (x *XSSChecker) fetch(URL string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (compatible; AssassinGo/0.1)")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body)
}
