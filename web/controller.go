package web

import (
	"../assassin"
	"../crawler"
	"../gatherer"
	"github.com/AmyangXYZ/sweetygo"
)

func home(ctx *sweetygo.Context) {
	ctx.Text(200, "biu")
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
	ret := map[string][]string{
		"emails":       emails,
		"fuzzableURLs": urls,
	}
	ctx.JSON(200, ret, "success")
}
