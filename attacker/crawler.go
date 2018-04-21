package attacker

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"../logger"
	"../util"
	"github.com/gorilla/websocket"
)

// Crawler crawls the website.
// WebSocket API.
type Crawler struct {
	mconn       *util.MuxConn
	host        string
	visitedURLs sync.Map
	emails      sync.Map
	// extractURLsRe finds the next urls to crawl.
	extractURLsRe *regexp.Regexp
	// extract emails.
	extractEmailsRe *regexp.Regexp
	// replaceRe relace params value in url.
	replaceGETValuesRe *regexp.Regexp
	maxDepth           int
	results            []string
}

// NewCrawler returns a new crawler.
func NewCrawler() *Crawler {
	return &Crawler{
		extractURLsRe:      regexp.MustCompile(`(?s)<a[ t]+.*?href="(.*?)".*?>`),
		replaceGETValuesRe: regexp.MustCompile(`(\?|\&)([^=]+)\=([^&]+)`),
	}
}

// Set implements Attacker interface.
// Params should be {conn *websocket.Conn, host string, depth int}
func (c *Crawler) Set(v ...interface{}) {
	c.mconn.Conn = v[0].(*websocket.Conn)
	c.host = "http://" + v[1].(string)
	c.maxDepth = v[2].(int)
}

// Report implements Attacker interface
func (c *Crawler) Report() map[string]interface{} {
	return map[string]interface{}{
		"fuzzableURLs": c.results,
	}
}

// Run implements Attacker interface.
func (c *Crawler) Run() {
	logger.Green.Println("Fuzzable URLs Crawling...")
	c.results = []string{}

	results := make(chan string)
	go c.Crawl(c.host, c.maxDepth, results)

	for url := range results {
		logger.Blue.Println(url)
		ret := map[string]string{
			"url": url,
		}
		c.mconn.Send(ret)
		c.results = append(c.results, url)
	}

	if len(c.results) == 0 {
		logger.Blue.Println("no fuzzable urls found")
	}
}

// Crawl crawls the target.
func (c *Crawler) Crawl(URL string, depth int, ret chan string) {
	defer close(ret)

	if depth <= 0 {
		return
	}

	// filter paramed url.
	tmpURL := c.replaceGETValuesRe.ReplaceAllString(URL, `$2`)

	// if url has been visited
	if _, ok := c.visitedURLs.Load(tmpURL); ok {
		return
	}
	c.visitedURLs.Store(tmpURL, true)

	if strings.Contains(URL, "?") {
		ret <- URL
	}

	nextURLsMap := c.fetch(URL)

	var nextURLs []string
	for _, nextURL := range nextURLsMap {
		nextURLs = append(nextURLs, nextURL)
	}

	results := make([]chan string, len(nextURLs))
	for i, next := range nextURLs {
		results[i] = make(chan string)
		go c.Crawl(next, depth-1, results[i])
	}

	for i := range results {
		for s := range results[i] {
			ret <- s
		}
	}

	return
}

// fetch the page and extract emails and next urls.
func (c *Crawler) fetch(URL string) map[string]string {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0")
	resp, err := client.Do(req)
	if err != nil {
		return map[string]string{}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	nextURLsMap := c.extractURLs(URL, string(body))
	return nextURLsMap
}

func (c *Crawler) extractURLs(URL, body string) map[string]string {
	extractedURLs := c.extractURLsRe.FindAllStringSubmatch(body, -1)
	u := ""
	URLs := make(map[string]string)
	baseURL, _ := url.Parse(URL)
	if extractedURLs != nil {
		for _, z := range extractedURLs {
			u = z[1]
			ur, err := url.Parse(z[1])
			if err == nil {
				if u == "/" {
					u = ""
				} else if ur.IsAbs() == true {
					continue
				} else if ur.IsAbs() == false {
					u = baseURL.ResolveReference(ur).String()
				} else if strings.HasPrefix(u, "//") {
					u = "http:" + u
				} else if strings.HasPrefix(u, "/") {
					u = c.host + u
				} else {
					u = URL + u
				}
				if strings.Contains(u, c.host) {
					URLs[c.replaceGETValuesRe.ReplaceAllString(u, `$2`)] = u
				}
			}
		}
	}
	return URLs
}
