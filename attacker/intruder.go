package attacker

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/AmyangXYZ/barbarian"

	"../logger"
	"../utils"
	"github.com/gorilla/websocket"
)

// Intruder intrudes the target.
// WebSocket API.
type Intruder struct {
	mconn       *utils.MuxConn
	target      string
	header      string
	intrudeType string
	re          *regexp.Regexp
	payloads    []string
	concurrency int
}

// NewIntruder returns a new intruder.
func NewIntruder() *Intruder {
	return &Intruder{
		mconn: &utils.MuxConn{},
		re:    regexp.MustCompile(`\$\$(.*?)\$\$`),
	}
}

// Set sets params for intruder.
// Params should be {conn *websocket.Conn, target, header, payload string, concurrency int}
func (i *Intruder) Set(v ...interface{}) {
	i.mconn.Conn = v[0].(*websocket.Conn)
	i.target = v[1].(string)
	i.header = v[2].(string)
	i.payloads = strings.Split(v[3].(string), "\n")
	i.concurrency = v[4].(int)
}

// Report implements Attacker interface.
func (i *Intruder) Report() map[string]interface{} {
	return nil
}

// Run implements Attacker interface.
func (i *Intruder) Run() {
	logger.Green.Println("Start Intruder...")
	bb := barbarian.New(i.attack, i.onResult, i.payloads, i.concurrency)
	var s utils.Signal
	go func() {
		i.mconn.Conn.ReadJSON(&s)
		if s.Stop == 1 {
			bb.Stop()
		}
	}()
	bb.Run()
}

func (i *Intruder) onResult(res interface{}) {
	ret := res.(map[string]string)
	logger.Blue.Println("Payload:", ret["payload"], "Status:",
		ret["resp_status"], "resp_len:", ret["len"])
	i.mconn.Send(res)
}

func (i *Intruder) attack(payload string) interface{} {
	resp := i.fetch(payload)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	res := map[string]string{
		"payload":     payload,
		"resp_status": strconv.Itoa(resp.StatusCode),
		"resp_len":    strconv.Itoa(len(string(body))),
	}
	return res
}

func (i *Intruder) fetch(payload string) *http.Response {
	client := &http.Client{}
	req := i.parse(payload)
	resp, err := client.Do(req)
	if err != nil {
		logger.Red.Println(err)
		return &http.Response{}
	}
	return resp
}

func (i *Intruder) parse(payload string) *http.Request {
	header := i.re.ReplaceAllString(i.header, payload)
	var x []string
	var data string
	// has body
	if strings.Contains(header, "\n\n") {
		x = strings.Split(header, "\n\n")
		x = append(x, header)
		data = x[1]
	} else {
		x = append(x, header)
	}

	hr := strings.Split(x[0], "\n")
	if len(hr) < 2 {
		logger.Red.Println("invalid header")
	}
	method := strings.Split(hr[0], " ")[0]
	path := strings.Split(hr[0], " ")[1]
	req, _ := http.NewRequest(method, "http://"+i.target+path, strings.NewReader(data))

	for _, row := range hr[1:] {
		hh := strings.Split(row, ": ")
		if len(hh) > 1 {
			k := hh[0]
			v := hh[1]
			req.Header.Add(k, v)
		}
	}
	return req
}
