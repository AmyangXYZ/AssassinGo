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
	mconn  *muxConn
	target string
	CMS    string
}

// NewCMSDetector returns a new CMS detector.
func NewCMSDetector() *CMSDetector {
	return &CMSDetector{}
}

// Set implements Gatherer interface.
// Params should be {conn *websocket.Conn, target string}
func (c *CMSDetector) Set(v ...interface{}) {
	c.mconn = &muxConn{conn: v[0].(*websocket.Conn)}
	c.target = v[1].(string)
}

// Report implements Gatherer interface
func (c *CMSDetector) Report() interface{} {
	return c.CMS
}

// Run impplements Gatherer interface.
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
	ret := map[string]string{
		"cms": c.CMS,
	}
	c.mconn.send(ret)
	logger.Green.Println("CMS Detected:", c.CMS)
}
