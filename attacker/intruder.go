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
type Intruder struct {
	target          string
	header          string
	intrudeType     string
	re              *regexp.Regexp
	payload         []string
	goroutinesCount int
}

// NewIntruder returns a new intruder.
func NewIntruder() *Intruder {
	return &Intruder{re: regexp.MustCompile(`\$\$(.*?)\$\$`)}
}

// Set sets params for intruder.
// Params should be {target, header, payload string, goroutinesCount int}
func (i *Intruder) Set(v ...interface{}) {
	i.target = v[0].(string)
	i.header = v[1].(string)
	i.payload = strings.Split(v[2].(string), ",")
	i.goroutinesCount = v[3].(int)
}

// Report implements attacker interface.
func (i *Intruder) Report() interface{} {
	return ""
}

// Run implements attacker interface.
func (i *Intruder) Run(conn *websocket.Conn) {
	logger.Green.Println("Start Intruder...")
	blockers := make(chan struct{}, i.goroutinesCount)
	for _, p := range i.payload {
		blockers <- struct{}{}
		go i.attack(conn, p, blockers)
	}

	for i := 0; i < cap(blockers); i++ {
		blockers <- struct{}{}
	}
}

func (i *Intruder) attack(conn *websocket.Conn, payload string, blocker chan struct{}) {
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
	conn.WriteJSON(ret)
}

func (i *Intruder) fetch(payload string) *http.Response {
	client := &http.Client{}
	req := i.parse(payload)
	resp, err := client.Do(req)
	if err != nil {
		logger.Red.Println(err)
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
