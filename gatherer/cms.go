package gatherer

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"../logger"
)

// CMSDetector detects CMS with whatcms.org api.
type CMSDetector struct {
	target string
	CMS    string
}

// NewCMSDetector returns a new CMS detector.
func NewCMSDetector(target string) *CMSDetector {
	return &CMSDetector{target: target}
}

// Report impletements Gatherer interface
func (c *CMSDetector) Report() interface{} {
	return c.CMS
}

// Run imppletements Gatherer interface.
func (c *CMSDetector) Run() {
	resp, err := http.Get("https://whatcms.org/?s=" + c.target)
	if err != nil {
		logger.Red.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	re := regexp.MustCompile(`class="nowrap" title="(.*?)">`)
	cms := re.FindAllStringSubmatch(string(body), -1)
	if len(cms) == 0 {
		logger.Green.Println("No CMS Detected")
		return
	}
	c.CMS = cms[0][1]
	logger.Green.Println("CMS Detected:", c.CMS)
}
