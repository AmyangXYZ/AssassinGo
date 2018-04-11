package seeker

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
	"github.com/gorilla/websocket"
)

// Run starts seek target using search engine.
func Run(q, se string, maxPage int, conn *websocket.Conn) {
	var err error
	if err != nil {

	}
	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	// c, err := chromedp.New(ctxt)
	c, err := chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)))

	urls := []string{}
	if se == "google" {
		err = c.Run(ctxt, searchGoogle(q, &urls, maxPage, conn))
	} else if se == "bing" {
		err = c.Run(ctxt, searchBing(q, &urls, maxPage, conn))
	}

	// fmt.Println(urls)
	// shutdown chrome
	err = c.Shutdown(ctxt)

	// wait for chrome to finish
	err = c.Wait()
}

func searchBing(q string, urls *[]string, page int, conn *websocket.Conn) chromedp.Tasks {
	var res []string
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.bing.com`),
		chromedp.Sleep(2 * time.Second),
		chromedp.SendKeys(`#sb_form_q`, q+"\n", chromedp.ByID),
		chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(c context.Context, e cdp.Executor) error {
			var resCount string
			chromedp.Text(`.sb_count`, &resCount, chromedp.ByQuery).Do(c, e)
			fmt.Println(resCount)
			n := strings.Replace(strings.Split(resCount, " ")[0], ",", "", -1)
			count, _ := strconv.Atoi(n)
			p := int(math.Floor(float64(count / 10)))
			if p < page {
				page = p
			}
			for i := 0; i <= page; i++ {
				chromedp.Sleep(1*time.Second).Do(c, e)
				chromedp.EvaluateAsDevTools(`
					var h2 = document.getElementsByTagName('h2');
					var urls = [];
					for(var i=0;i<h2.length-2;i++){
						var a = h2[i].getElementsByTagName('a');
						urls.push(a[0].href);
					}
					urls`, &res).Do(c, e)

				ret := map[string][]string{
					"urls": res,
				}
				conn.WriteJSON(ret)

				*urls = append(*urls, res...)
				if i != page {
					chromedp.Click(`//*[@title="Next page"]`, chromedp.BySearch).Do(c, e)
				}
			}
			return nil
		}),
	}
}

func searchGoogle(q string, urls *[]string, page int, conn *websocket.Conn) chromedp.Tasks {
	res := []string{}
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.google.com`),
		chromedp.Sleep(2 * time.Second),
		chromedp.SendKeys(`#lst-ib`, q+"\n", chromedp.ByID),
		chromedp.WaitVisible(`#res`, chromedp.ByID),
		chromedp.ActionFunc(func(c context.Context, e cdp.Executor) error {
			var resCount string
			var ss string
			chromedp.Text(`#resultStats`, &resCount, chromedp.ByID).Do(c, e)
			s := strings.Split(resCount, " ")
			if s[1] == "results" {
				ss = s[0]
			} else {
				ss = s[1]
			}
			n := strings.Replace(ss, ",", "", -1)
			count, _ := strconv.Atoi(n)
			p := int(math.Floor(float64(count / 10)))
			if p < page {
				page = p
			}
			for i := 0; i <= page; i++ {
				chromedp.Sleep(2*time.Second).Do(c, e)
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
					urls`, &res).Do(c, e)

				ret := map[string][]string{
					"urls": res,
				}
				conn.WriteJSON(ret)

				*urls = append(*urls, res...)
				if i != page {
					chromedp.Click(`#pnnext`, chromedp.NodeVisible).Do(c, e)
				}
			}
			return nil
		}),
	}
}
