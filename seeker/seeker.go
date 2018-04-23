package seeker

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"../logger"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
	"github.com/chromedp/chromedp/runner"
	"github.com/gorilla/websocket"
)

// Seeker seeks targets with search engine.
type Seeker struct {
	conn    *websocket.Conn
	query   string
	se      string
	maxPage int
	Results []string
}

// NewSeeker returns a new Seeker.
func NewSeeker(q, se string, maxPage int) *Seeker {
	return &Seeker{
		query:   q,
		se:      se,
		maxPage: maxPage,
	}
}

// Set params.
// Params should be {conn *websocket.Conn, query ,se stirng, maxPage int}
func (s *Seeker) Set(v ...interface{}) {
	s.conn = v[0].(*websocket.Conn)
	s.query = v[1].(string)
	s.se = v[2].(string)
	s.maxPage = v[3].(int)
}

// Run starts seeker.
func (s *Seeker) Run() {
	logger.Green.Println("Seeking Targets...")
	logger.Blue.Println("Search Engine:", s.se)
	logger.Blue.Println("Keyword:", s.query)
	logger.Blue.Println("Max Page:", s.maxPage)
	var err error
	if err != nil {

	}
	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	options := chromedp.WithRunnerOptions(
		runner.Flag("no-first-run", true),
		runner.Flag("no-sandbox", true),
		runner.Flag("disable-gpu", true),
	)
	// create chrome instance
	// c, err := chromedp.New(ctxt, options)
	c, err := chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)), options)
	if err != nil {
		logger.Red.Println(err)
		return
	}

	if s.se == "google" {
		err = c.Run(ctxt, s.searchGoogle())
	} else if s.se == "bing" {
		err = c.Run(ctxt, s.searchBing())
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)

	// wait for chrome to finish
	err = c.Wait()
}

func (s *Seeker) searchBing() chromedp.Tasks {
	var urls []string
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.bing.com`),
		chromedp.Sleep(2 * time.Second),
		chromedp.SendKeys(`#sb_form_q`, s.query+"\n", chromedp.ByID),
		chromedp.WaitVisible(`.sb_count`, chromedp.ByQuery),
		chromedp.ActionFunc(func(c context.Context, e cdp.Executor) error {
			var resCount string
			chromedp.Text(`.sb_count`, &resCount, chromedp.ByQuery).Do(c, e)
			n := strings.Replace(strings.Split(resCount, " ")[0], ",", "", -1)
			count, _ := strconv.Atoi(n)
			p := int(math.Floor(float64(count / 10)))
			if p < s.maxPage {
				s.maxPage = p
			}
			for i := 0; i <= s.maxPage; i++ {
				chromedp.Sleep(2*time.Second).Do(c, e)
				chromedp.EvaluateAsDevTools(`
					var h2 = document.getElementsByTagName('h2');
					var urls = [];
					for(var i=0;i<h2.length-2;i++){
						var a = h2[i].getElementsByTagName('a');
						urls.push(a[0].href);
					}
					urls`, &urls).Do(c, e)
				for _, u := range urls {
					logger.Blue.Println(u)
				}
				ret := map[string][]string{
					"urls": urls,
				}
				s.conn.WriteJSON(ret)

				s.Results = append(s.Results, urls...)
				if i != s.maxPage {
					chromedp.Click(`//*[@title="Next page"]`, chromedp.BySearch).Do(c, e)
				}
			}
			return nil
		}),
	}
}

func (s *Seeker) searchGoogle() chromedp.Tasks {
	urls := []string{}
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.google.com`),
		chromedp.SendKeys(`#lst-ib`, s.query+"\n", chromedp.ByID),
		chromedp.WaitVisible(`#res`, chromedp.ByID),
		chromedp.ActionFunc(func(c context.Context, e cdp.Executor) error {
			var resCount string
			var xx string
			chromedp.Text(`#resultStats`, &resCount, chromedp.ByID).Do(c, e)
			x := strings.Split(resCount, " ")
			if x[1] == "results" {
				xx = x[0]
			} else {
				xx = x[1]
			}
			n := strings.Replace(xx, ",", "", -1)
			count, _ := strconv.Atoi(n)
			p := int(math.Floor(float64(count / 10)))
			if p < s.maxPage {
				s.maxPage = p
			}
			for i := 0; i <= s.maxPage; i++ {
				chromedp.Sleep(1*time.Second).Do(c, e)
				chromedp.EvaluateAsDevTools(`
					var h3 = document.getElementsByTagName('h3');
					var c = h3.length;
					if(h3.length==11){
						c=10;
					}
					var urls = [];
					for(var i=0;i<c;i++){
						var a = h3[i].getElementsByTagName('a');
						urls.push(a[0].href);
					}
					urls`, &urls).Do(c, e)
				for _, u := range urls {
					logger.Blue.Println(u)
				}
				ret := map[string][]string{
					"urls": urls,
				}
				s.conn.WriteJSON(ret)

				s.Results = append(s.Results, urls...)
				if i != s.maxPage {
					chromedp.Click(`#pnnext`, chromedp.NodeVisible).Do(c, e)
				}
			}
			return nil
		}),
	}
}
