package seeker

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"assassingo/logger"
	"assassingo/utils"

	"github.com/chromedp/chromedp"
	"github.com/gorilla/websocket"
)

// Seeker seeks targets with search engine.
type Seeker struct {
	mconn   *utils.MuxConn
	query   string
	se      string
	maxPage int
	Results []string
}

// NewSeeker returns a new Seeker.
func NewSeeker(q, se string, maxPage int) *Seeker {
	return &Seeker{
		mconn:   &utils.MuxConn{},
		query:   q,
		se:      se,
		maxPage: maxPage,
	}
}

// Set params.
// Params should be {conn *websocket.Conn, query, se string, maxPage int}
func (s *Seeker) Set(v ...interface{}) {
	s.mconn.Conn = v[0].(*websocket.Conn)
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

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	var err error
	if s.se == "google" {
		err = chromedp.Run(ctx, s.searchGoogle())
	} else if s.se == "bing" {
		err = chromedp.Run(ctx, s.searchBing())
	}

	if err != nil {
		logger.Red.Println(err)
	}
}

func (s *Seeker) searchBing() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.bing.com`),
		chromedp.Sleep(2 * time.Second),
		chromedp.SendKeys(`#sb_form_q`, s.query+"\n"),
		chromedp.WaitVisible(`.sb_count`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var resCount string
			if err := chromedp.Text(`.sb_count`, &resCount).Do(ctx); err != nil {
				return err
			}
			n := strings.Replace(strings.Split(resCount, " ")[0], ",", "", -1)
			count, _ := strconv.Atoi(n)
			p := int(math.Floor(float64(count / 10)))
			if p < s.maxPage {
				s.maxPage = p
			}
			for i := 0; i <= s.maxPage; i++ {
				if err := chromedp.Sleep(2 * time.Second).Do(ctx); err != nil {
					return err
				}
				var urls []string
				if err := chromedp.Evaluate(`
					Array.from(document.getElementsByTagName('h2'))
						.slice(0, -2)
						.map(h2 => h2.getElementsByTagName('a')[0].href)
				`, &urls).Do(ctx); err != nil {
					return err
				}
				for _, u := range urls {
					logger.Blue.Println(u)
				}
				s.mconn.Send(map[string][]string{"urls": urls})
				s.Results = append(s.Results, urls...)
				if i != s.maxPage {
					if err := chromedp.Click(`//*[@title="Next page"]`, chromedp.BySearch).Do(ctx); err != nil {
						return err
					}
				}
			}
			return nil
		}),
	}
}

func (s *Seeker) searchGoogle() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.google.com`),
		chromedp.SendKeys(`input[name=q]`, s.query+"\n"),
		chromedp.WaitVisible(`#search`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var resCount string
			if err := chromedp.Text(`#result-stats`, &resCount).Do(ctx); err != nil {
				return err
			}
			parts := strings.Fields(resCount)
			var countStr string
			if len(parts) >= 2 && parts[1] == "results" {
				countStr = parts[0]
			} else if len(parts) >= 3 {
				countStr = parts[1]
			}
			count, _ := strconv.Atoi(strings.ReplaceAll(countStr, ",", ""))
			p := int(math.Floor(float64(count / 10)))
			if p < s.maxPage {
				s.maxPage = p
			}
			for i := 0; i <= s.maxPage; i++ {
				if err := chromedp.Sleep(1 * time.Second).Do(ctx); err != nil {
					return err
				}
				var urls []string
				if err := chromedp.Evaluate(`
					Array.from(document.getElementsByTagName('h3'))
						.slice(0, 10)
						.map(h3 => h3.getElementsByTagName('a')[0].href)
				`, &urls).Do(ctx); err != nil {
					return err
				}
				for _, u := range urls {
					logger.Blue.Println(u)
				}
				s.mconn.Send(map[string][]string{"urls": urls})
				s.Results = append(s.Results, urls...)
				if i != s.maxPage {
					if err := chromedp.Click(`#pnnext`, chromedp.NodeVisible).Do(ctx); err != nil {
						return err
					}
				}
			}
			return nil
		}),
	}
}
