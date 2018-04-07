package crawler

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"../logger"
	"github.com/gorilla/websocket"
)

// Crawler crawls the website.
type Crawler struct {
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
}

// NewCrawler returns a new crawler.
func NewCrawler(host string, depth int) Crawler {
	return Crawler{
		host:               "http://" + host,
		extractURLsRe:      regexp.MustCompile(`(?s)<a[ t]+.*?href="(.*?)".*?>`),
		replaceGETValuesRe: regexp.MustCompile(`(\?|\&)([^=]+)\=([^&]+)`),
		maxDepth:           depth,
	}
}

// Run begins crawling and return fuzzable urls.
func (c *Crawler) Run(conn *websocket.Conn) []string {
	logger.Green.Println("Fuzzable URLs Crawling...")
	var fuzzableURLs []string
	results := make(chan string)

	go c.Crawl(c.host, c.maxDepth, results)
	for url := range results {
		logger.Blue.Println(url)
		conn.WriteJSON(url)
		fuzzableURLs = append(fuzzableURLs, url)
	}

	if len(fuzzableURLs) == 0 {
		logger.Blue.Println("no fuzzable urls found")
	}

	return fuzzableURLs
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
	req.Header.Set("user-agent", "Mozilla/5.0 (compatible; AssassinGo/0.1)")
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
