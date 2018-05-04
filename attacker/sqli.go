package attacker

import (
	"io/ioutil"
	"net/http"

	"github.com/AmyangXYZ/barbarian"

	"../logger"
	"../utils"
	"github.com/gorilla/websocket"
)

// BasicSQLi checks basic sqli vuls.
// WebSocket API.
type BasicSQLi struct {
	mconn         *utils.MuxConn
	fuzzableURLs  []string
	payload0      string
	payload1      string
	InjectableURL []string
}

// NewBasicSQLi returns a new basicSQli Attacker.
func NewBasicSQLi() *BasicSQLi {
	return &BasicSQLi{
		mconn:    &utils.MuxConn{},
		payload0: "/**/%26%261%3d2%23",
		payload1: "/**/%26%261%3d1%23",
	}
}

// Set implements Attacker interface.
// Params should be {conn *websocket.Conn, fuzzableURLs []string}
func (bs *BasicSQLi) Set(v ...interface{}) {
	bs.mconn.Conn = v[0].(*websocket.Conn)
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
	bb := barbarian.New(bs.check, bs.onResult, bs.fuzzableURLs, 20)
	bb.Run()
	if len(bs.InjectableURL) == 0 {
		logger.Blue.Println("no sqli vuls found")
	}
}

func (bs *BasicSQLi) onResult(res interface{}) {
	logger.Blue.Println(res)
	ret := map[string]string{
		"sqli_url": res.(string),
	}
	bs.mconn.Send(ret)
	bs.InjectableURL = append(bs.InjectableURL, res.(string))
}

func (bs *BasicSQLi) check(URL string) interface{} {
	body0 := bs.fetch(URL + bs.payload0)
	body1 := bs.fetch(URL + bs.payload1)
	if len(body0) != len(body1) {
		return URL

	}
	return nil
}

func (bs *BasicSQLi) fetch(URL string) string {
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
