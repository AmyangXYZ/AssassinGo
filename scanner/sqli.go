package scanner

import (
	"io/ioutil"
	"net/http"

	"../logger"
)

// BasicSQLi checks basic sqli vuls.
type BasicSQLi struct {
	payload0      string
	payload1      string
	InjectableURL []string
}

// NewBasicSQli returns a new basicSQli scanner.
func NewBasicSQli() *BasicSQLi {
	return &BasicSQLi{
		payload0: "/**/%26%261%3d2%23",
		payload1: "/**/%26%261%3d1%23",
	}
}

// Run impletements Scanner interface.
func (bs *BasicSQLi) Run(fuzzableURLs []string) {
	logger.Green.Println("Basic SQLi Checking...")

	blockers := make(chan bool, len(fuzzableURLs))
	for _, URL := range fuzzableURLs {
		blockers <- true
		go bs.check(URL, blockers)
	}

	// Wait for all goroutines to finish.
	for i := 0; i < cap(blockers); i++ {
		blockers <- true
	}
	if len(bs.InjectableURL) == 0 {
		logger.Blue.Println("no sqli vuls found")
	}
}

func (bs *BasicSQLi) check(URL string, blocker chan bool) {
	defer func() { <-blocker }()
	body0 := bs.fetch(URL + bs.payload0)
	body1 := bs.fetch(URL + bs.payload1)
	if len(body0) != len(body1) {
		logger.Blue.Println(URL)
		bs.InjectableURL = append(bs.InjectableURL, URL)
	}
}

func (bs *BasicSQLi) fetch(URL string) string {
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
