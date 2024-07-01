package gatherer

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"logger"
)

// Honeypot detects honeypot score of the target.
// AJAX API.
type Honeypot struct {
	target string
	Score  string
}

// NewHoneypotDetecter returns a new honeypot detecter.
func NewHoneypotDetecter() *Honeypot {
	return &Honeypot{}
}

// Set implements Gatherer interface.
// Params should be {targetIP string}
func (h *Honeypot) Set(v ...interface{}) {
	h.target = v[0].(string)
}

// Report implements Gatherer interface.
func (h *Honeypot) Report() map[string]interface{} {
	return map[string]interface{}{
		"score": h.Score,
	}
}

// Run implements Gatherer interface.
func (h *Honeypot) Run() {
	h.honeypotDetect()
	logger.Green.Println("Honeypot Score:", h.Score)
}

func (h *Honeypot) honeypotDetect() {
	targetIP, err := net.ResolveIPAddr("ip", h.target)
	if err != nil {
		logger.Red.Println(err)
		return
	}
	u := fmt.Sprintf("https://api.shodan.io/labs/honeyscore/%s?key=6Whsxxn9Ajjuc5nHB0CDWTdPIOKSJ0zy", targetIP)
	resp, err := http.Get(u)
	if err != nil {
		logger.Red.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	h.Score = string(body)
}
