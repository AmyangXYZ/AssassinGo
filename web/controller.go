package web

import (
	"net/http"

	"../assassin"
	"../crawler"
	"../gatherer"
	"../scanner"
	"github.com/AmyangXYZ/sweetygo"
)

func home(ctx *sweetygo.Context) {
	ctx.Render(200, "home")
}

func seek(ctx *sweetygo.Context) {
	ctx.Render(200, "seek")
}

func shadow(ctx *sweetygo.Context) {
	ctx.Render(200, "shadow")
}

func attack(ctx *sweetygo.Context) {
	ctx.Render(200, "attack")
}

func assassinate(ctx *sweetygo.Context) {
	ctx.Render(200, "assassinate")
}

func static(ctx *sweetygo.Context) {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("/home/amyang/Projects/AssassinGo/web/static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
}

func newAssassin(ctx *sweetygo.Context) {
	target := ctx.Param("target")
	a = assassin.New(target)
	ctx.Text(200, target)
}

func basicInfo(ctx *sweetygo.Context) {
	B := gatherer.NewBasicInfo(a.Target)
	B.Run()
	bi := B.Report().([]string)
	ctx.JSON(200, bi, "success")
}

func cmsDetect(ctx *sweetygo.Context) {
	C := gatherer.NewCMSDetector(a.Target)
	C.Run()
	cms := C.Report().(string)
	ctx.JSON(200, map[string]string{"cms": cms}, "success")
}

func portScan(ctx *sweetygo.Context) {
	P := gatherer.NewPortScanner(a.Target)
	P.Run()
	ports := P.Report().([]string)
	ctx.JSON(200, ports, "success")
}

func crawl(ctx *sweetygo.Context) {
	C := crawler.NewCrawler(a.Target, 4)
	emails, urls := C.Run()
	a.FuzzableURLs = urls
	ret := map[string][]string{
		"emails":       emails,
		"fuzzableURLs": urls,
	}
	ctx.JSON(200, ret, "success")
}

func checkSQLi(ctx *sweetygo.Context) {
	S := scanner.NewBasicSQLi()
	S.Run(a.FuzzableURLs)
	urls := S.Report().([]string)
	ctx.JSON(200, urls, "success")
}

func checkXSS(ctx *sweetygo.Context) {
	X := scanner.NewXSSChecker()
	X.Run(a.FuzzableURLs)
	urls := X.Report().([]string)
	ctx.JSON(200, urls, "success")
}
