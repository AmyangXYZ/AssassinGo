package gatherer

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"../logger"
	"github.com/gorilla/websocket"
)

// CMSDetector detects CMS with whatcms.org api.
type CMSDetector struct {
	target string
	CMS    string
}

// NewCMSDetector returns a new CMS detector.
func NewCMSDetector() *CMSDetector {
	return &CMSDetector{}
}

// Set implements Gatherer interface.
// Params should be {target string}
func (c *CMSDetector) Set(v ...interface{}) {
	c.target = v[0].(string)
}

// Report implements Gatherer interface
func (c *CMSDetector) Report() interface{} {
	return c.CMS
}

// Run impplements Gatherer interface.
func (c *CMSDetector) Run(conn *websocket.Conn) {
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
	ret := map[string]string{
		"cms": c.CMS,
	}
	conn.WriteJSON(ret)
	logger.Green.Println("CMS Detected:", c.CMS)
}
