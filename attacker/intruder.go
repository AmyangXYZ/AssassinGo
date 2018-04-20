package attacker

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"../logger"
	"github.com/gorilla/websocket"
)

// Intruder intrudes the target.
// WebSocket API.
type Intruder struct {
	mconn           *muxConn
	target          string
	header          string
	intrudeType     string
	re              *regexp.Regexp
	payloads        []string
	goroutinesCount int
}

// NewIntruder returns a new intruder.
func NewIntruder() *Intruder {
	return &Intruder{re: regexp.MustCompile(`\$\$(.*?)\$\$`)}
}

// Set sets params for intruder.
// Params should be {conn *websocket.Conn, target, header, payload string, goroutinesCount int}
func (i *Intruder) Set(v ...interface{}) {
	i.mconn = &muxConn{conn: v[0].(*websocket.Conn)}
	i.target = v[1].(string)
	i.header = v[2].(string)
	i.payloads = strings.Split(v[3].(string), "\n")
	i.goroutinesCount = v[4].(int)
}

// Report implements Attacker interface.
func (i *Intruder) Report() map[string]interface{} {
	return nil
}

// Run implements Attacker interface.
func (i *Intruder) Run() {
	logger.Green.Println("Start Intruder...")
	var s signal
	stop := make(chan struct{}, 0)
	blockers := make(chan struct{}, i.goroutinesCount)
	go func() {
		i.mconn.conn.ReadJSON(&s)
		if s.Stop == 1 {
			stop <- struct{}{}
		}
	}()

loop:
	for _, p := range i.payloads {
		select {
		default:
			blockers <- struct{}{}
			go i.attack(p, blockers)
		case <-stop:
			break loop
		}
	}

	for i := 0; i < cap(blockers); i++ {
		blockers <- struct{}{}
	}
}

func (i *Intruder) attack(payload string, blocker chan struct{}) {
	defer func() { <-blocker }()
	resp := i.fetch(payload)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	logger.Blue.Println("Payload:", payload, "Status:", resp.StatusCode, "len:", len(string(body)))
	ret := map[string]string{
		"payload":     payload,
		"resp_status": strconv.Itoa(resp.StatusCode),
		"resp_len":    strconv.Itoa(len(string(body))),
	}
	i.mconn.send(ret)
}

func (i *Intruder) fetch(payload string) *http.Response {
	client := &http.Client{}
	req := i.parse(payload)
	resp, err := client.Do(req)
	if err != nil {
		logger.Red.Println(err)
		return nil
	}
	return resp
}

func (i *Intruder) parse(payload string) *http.Request {
	header := i.re.ReplaceAllString(i.header, payload)
	x := strings.Split(header, "\n\n")
	data := x[1]
	hr := strings.Split(x[0], "\n")

	method := strings.Split(hr[0], " ")[0]
	path := strings.Split(hr[0], " ")[1]
	req, _ := http.NewRequest(method, "http://"+i.target+":8888"+path, strings.NewReader(data))

	for _, row := range hr[1:] {
		hh := strings.Split(row, ": ")
		k := hh[0]
		v := hh[1]
		req.Header.Add(k, v)
	}
	return req
}
