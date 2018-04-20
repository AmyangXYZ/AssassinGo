package gatherer

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../logger"
	"github.com/gorilla/websocket"
)

// DirBruter brute force the dir.
// WebSocket API.
type DirBruter struct {
	mconn           *muxConn
	target          string
	payloads        []string
	goroutinesCount int
}

// NewDirBruter returns a new dirbruter.
func NewDirBruter() *DirBruter {
	return &DirBruter{payloads: readPayloadsFromFile("./gatherer/dict/dir-php.txt")}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target, dic string, goroutinesCount int}
func (d *DirBruter) Set(v ...interface{}) {
	d.mconn = &muxConn{conn: v[0].(*websocket.Conn)}
	d.target = v[1].(string)
	d.goroutinesCount = v[2].(int)
}

// Report implements Gatherer interface.
func (d *DirBruter) Report() map[string]interface{} {
	return nil
}

// Run implements Gatherer interface,
func (d *DirBruter) Run() {
	logger.Green.Println("Brute Force Dir")
	var s signal
	stop := make(chan struct{}, 0)
	blockers := make(chan struct{}, d.goroutinesCount)
	go func() {
		d.mconn.conn.ReadJSON(&s)
		if s.Stop == 1 {
			stop <- struct{}{}
		}
	}()

loop:
	for _, p := range d.payloads {
		select {
		default:
			blockers <- struct{}{}
			go d.fetch(p, blockers)
		case <-stop:
			break loop
		}
	}

	for i := 0; i < cap(blockers); i++ {
		blockers <- struct{}{}
	}
}

func (d *DirBruter) fetch(path string, blocker chan struct{}) {
	defer func() { <-blocker }()
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", "http://"+d.target+path, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode == 404 {
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	logger.Blue.Println("Path:", path, "Status:", resp.StatusCode, "len:", len(string(body)))
	ret := map[string]string{
		"path":   path,
		"status": strconv.Itoa(resp.StatusCode),
		"len":    strconv.Itoa(len(string(body))),
	}
	d.mconn.send(ret)
}

func readPayloadsFromFile(file string) []string {
	buf, _ := ioutil.ReadFile(file)
	p := strings.Split(string(buf), "\n")
	return p
}
