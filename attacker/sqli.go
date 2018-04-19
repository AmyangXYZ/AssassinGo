package attacker

import (
	"io/ioutil"
	"net/http"

	"../logger"
	"github.com/gorilla/websocket"
)

// BasicSQLi checks basic sqli vuls.
// WebSocket API.
type BasicSQLi struct {
	mconn         *muxConn
	fuzzableURLs  []string
	payload0      string
	payload1      string
	InjectableURL []string
}

// NewBasicSQLi returns a new basicSQli Attacker.
func NewBasicSQLi() *BasicSQLi {
	return &BasicSQLi{
		payload0: "/**/%26%261%3d2%23",
		payload1: "/**/%26%261%3d1%23",
	}
}

// Set implements Attacker interface.
// Params should be {conn *websocket.Conn, fuzzableURLs []string}
func (bs *BasicSQLi) Set(v ...interface{}) {
	bs.mconn = &muxConn{conn: v[0].(*websocket.Conn)}
	bs.fuzzableURLs = v[1].([]string)
}

// Report implements Attacker interface.
func (bs *BasicSQLi) Report() map[string]interface{} {
	return map[string]interface{}{
		"sqli_urls": bs.InjectableURL,
	}
}

// Run implements Attacker interface.
func (bs *BasicSQLi) Run() {
	logger.Green.Println("Basic SQLi Checking...")
	bs.InjectableURL = []string{}

	blockers := make(chan bool, len(bs.fuzzableURLs))
	for _, URL := range bs.fuzzableURLs {
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
		ret := map[string]string{
			"sqli_url": URL,
		}
		bs.mconn.send(ret)
		bs.InjectableURL = append(bs.InjectableURL, URL)
	}
}

func (bs *BasicSQLi) fetch(URL string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; AssassinGo/0.1)")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body)
}
