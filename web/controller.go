package web

import (
	"net/http"
	"strings"

	"../assassin"
	"../crawler"
	"../gatherer"
	"../poc"
	"../scanner"
	"github.com/AmyangXYZ/sweetygo"
	"github.com/gorilla/websocket"
)

func index(ctx *sweetygo.Context) {
	ctx.Render(200, "index")
}

func static(ctx *sweetygo.Context) {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./web/static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
}

func newAssassin(ctx *sweetygo.Context) {
	target := ctx.Param("target")
	a = assassin.New(target)
	ret := map[string]string{
		"target": target,
	}
	ctx.JSON(201, ret, "success")
}

func basicInfo(ctx *sweetygo.Context) {
	B := gatherer.NewBasicInfo(a.Target)
	B.Run()
	bi := B.Report().([]string)
	ret := map[string]string{
		"ip":        bi[0],
		"webServer": bi[1],
	}
	ctx.JSON(200, ret, "success")
}

func cmsDetect(ctx *sweetygo.Context) {
	C := gatherer.NewCMSDetector(a.Target)
	C.Run()
	cms := C.Report().(string)
	ret := map[string]string{
		"cms": cms,
	}
	ctx.JSON(200, ret, "success")
}

func portScan(ctx *sweetygo.Context) {
	P := gatherer.NewPortScanner(a.Target)
	P.Run()
	ports := P.Report().([]string)
	ret := map[string][]string{
		"ports": ports,
	}
	ctx.JSON(200, ret, "success")
}

func crawl(ctx *sweetygo.Context) {
	C := crawler.NewCrawler(a.Target, 4)
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	urls := C.Run(conn)
	a.FuzzableURLs = urls
	conn.Close()
}

func checkSQLi(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	S := scanner.NewBasicSQLi()
	S.Run(a.FuzzableURLs, conn)
	conn.Close()
}

func checkXSS(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	X := scanner.NewXSSChecker()
	X.Run(a.FuzzableURLs, conn)
	conn.Close()
}

func intrude(ctx *sweetygo.Context) {

}

// POST -d "targets=t1,t2,t3..."
// batch scan is only for poc.
func setTargets(ctx *sweetygo.Context) {
	params := ctx.Params()
	ts := params["targets"][0]
	targets := strings.Split(ts, ",")
	for _, t := range targets {
		ateam = append(ateam, assassin.New(t))
	}

	ret := map[string][]string{
		"targets": targets,
	}
	ctx.JSON(201, ret, "success")
}

func getPOCs(ctx *sweetygo.Context) {
	var pocList []string
	for pocNames := range poc.POCMap {
		pocList = append(pocList, pocNames)
	}

	ret := map[string][]string{
		"poclist": pocList,
	}
	ctx.JSON(200, ret, "success")
}

// POST -d "poc=xxx"
func setPOC(ctx *sweetygo.Context) {
	pocName := ctx.Param("poc")
	for _, aa := range ateam {
		aa.POC = poc.POCMap[pocName]
	}

	ret := map[string]string{
		"poc": pocName,
	}
	ctx.JSON(201, ret, "success")
}

func runPOC(ctx *sweetygo.Context) {
	concurrency := 2
	blockers := make(chan struct{}, concurrency)
	var existedList []string

	for _, aa := range ateam {
		blockers <- struct{}{}
		go func(a *assassin.Assassin, blocker chan struct{}) {
			defer func() { <-blocker }()
			a.POC.Run(a.Target)
			if result := a.POC.Report().(string); result == "true" {
				existedList = append(existedList, a.Target)
			}
		}(aa, blockers)
	}
	for i := 0; i < cap(blockers); i++ {
		blockers <- struct{}{}
	}

	ret := map[string][]string{
		"existed": existedList,
	}
	ctx.JSON(200, ret, "success")
}
