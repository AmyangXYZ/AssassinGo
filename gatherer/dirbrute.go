package gatherer

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/AmyangXYZ/barbarian"

	"assassingo/logger"
	"assassingo/utils"

	"github.com/gorilla/websocket"
)

// DirBruter brute force the dir.
// WebSocket API.
type DirBruter struct {
	mconn       *utils.MuxConn
	target      string
	payloads    []string
	concurrency int
}

// NewDirBruter returns a new dirbruter.
func NewDirBruter() *DirBruter {
	return &DirBruter{
		mconn:    &utils.MuxConn{},
		payloads: utils.ReadFile("/dict/dir-php.txt"),
	}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target, dic string, concurrency int}
func (d *DirBruter) Set(v ...interface{}) {
	d.mconn.Conn = v[0].(*websocket.Conn)
	d.target = v[1].(string)
	d.concurrency = v[2].(int)
}

// Report implements Gatherer interface.
func (d *DirBruter) Report() map[string]interface{} {
	return nil
}

// Run implements Gatherer interface,
func (d *DirBruter) Run() {
	logger.Green.Println("Brute Force Dir")
	var s utils.Signal
	bb := barbarian.New(d.fetch, d.onResult, d.payloads, d.concurrency)
	go func() {
		d.mconn.Conn.ReadJSON(&s)
		if s.Stop == 1 {
			bb.Stop()
		}
	}()
	bb.Run()
}

func (d *DirBruter) onResult(res interface{}) {
	ret := res.(map[string]string)
	logger.Blue.Println("Path:", ret["path"],
		"Status:", ret["status"], "len:", ret["len"])
	d.mconn.Send(ret)
}

func (d *DirBruter) fetch(path string) interface{} {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", "http://"+d.target+path, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0")
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == 404 {
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	ret := map[string]string{
		"path":   path,
		"status": strconv.Itoa(resp.StatusCode),
		"len":    strconv.Itoa(len(string(body))),
	}
	return ret
}
