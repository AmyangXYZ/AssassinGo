package web

import (
	"net/http"
	"strings"

	"../assassin"
	"../logger"
	"../poc"
	"../seeker"
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
	if strings.Contains(target, ",") == false {
		a = assassin.New(target)
		ret := map[string]string{
			"target": target,
		}
		ctx.JSON(201, ret, "success")
		return
	}
	targets := strings.Split(target, ",")
	for _, t := range targets {
		ateam = append(ateam, assassin.New(t))
	}

	ret := map[string][]string{
		"targets": targets,
	}
	ctx.JSON(201, ret, "success")
}

func basicInfo(ctx *sweetygo.Context) {
	a.Gatherers["basicInfo"].Set(a.Target)
	a.Gatherers["basicInfo"].Run()
	ret := a.Gatherers["basicInfo"].Report()

	ctx.JSON(200, ret, "success")
}

func cmsDetect(ctx *sweetygo.Context) {
	a.Gatherers["cms"].Set(a.Target)
	a.Gatherers["cms"].Run()
	ret := a.Gatherers["cms"].Report()
	ctx.JSON(200, ret, "success")
}

func whois(ctx *sweetygo.Context) {
	a.Gatherers["whois"].Set(a.Target)
	a.Gatherers["whois"].Run()
	ret := a.Gatherers["whois"].Report()
	ctx.JSON(200, ret, "success")
}

func honeypot(ctx *sweetygo.Context) {
	a.Gatherers["honeypot"].Set(a.Target)
	a.Gatherers["honeypot"].Run()
	ret := a.Gatherers["honeypot"].Report()
	ctx.JSON(200, ret, "success")
}

func tracert(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["tracert"].Set(conn, a.Target)
	a.Gatherers["tracert"].Run()
	conn.Close()
}

func portScan(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["port"].Set(conn, a.Target, "tcp")
	a.Gatherers["port"].Run()
	conn.Close()
}

type dirbMsg struct {
	// Payload   string `json:"payload"`
	GortCount int `json:"gort_count"`
}

func dirBrute(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := dirbMsg{}
	conn.ReadJSON(&m)
	a.Gatherers["dirb"].Set(conn, a.Target, m.GortCount)
	a.Gatherers["dirb"].Run()
	conn.Close()
}

func crawl(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["crawl"].Set(conn, a.Target, 4)
	a.Attackers["crawl"].Run()
	a.FuzzableURLs = a.Attackers["crawl"].Report()["fuzzableURLs"].([]string)
}

func checkSQLi(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["sqli"].Set(conn, a.FuzzableURLs)
	a.Attackers["sqli"].Run()
	conn.Close()
}

func checkXSS(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["xss"].Set(conn, a.FuzzableURLs)
	a.Attackers["xss"].Run()
	conn.Close()
}

type intruderMsg struct {
	Header    string `json:"header"`
	Payload   string `json:"payload"`
	GortCount int    `json:"gort_count"`
}

func intrude(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := intruderMsg{}
	conn.ReadJSON(&m)
	a.Attackers["intruder"].Set(conn, a.Target, m.Header, m.Payload, m.GortCount)
	a.Attackers["intruder"].Run()
	conn.Close()
}

type seekerMsg struct {
	Query   string `json:"query"`
	SE      string `json:"se"`
	MaxPage int    `json:"max_page"`
}

func seek(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := seekerMsg{}
	conn.ReadJSON(&m)
	S := seeker.NewSeeker(m.Query, m.SE, m.MaxPage)
	S.Run(conn)
	conn.Close()
}

// POST -d "poc=xxx"
func setPoC(ctx *sweetygo.Context) {
	pocName := ctx.Param("poc")
	if len(ateam) > 0 {
		for _, aa := range ateam {
			aa.PoC = poc.PoCMap[pocName]
		}
	} else {
		a.PoC = poc.PoCMap[pocName]
	}
	logger.Green.Println("PoC Set:", pocName)
	ret := map[string]string{
		"poc": pocName,
	}
	ctx.JSON(201, ret, "success")
}

func runPoC(ctx *sweetygo.Context) {
	if len(ateam) > 0 {
		conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
		concurrency := 10
		blockers := make(chan struct{}, concurrency)
		var exploitableList []string
		for _, aa := range ateam {
			blockers <- struct{}{}
			go func(a *assassin.Assassin, blocker chan struct{}) {
				defer func() { <-blocker }()
				a.PoC.Set(conn, a.Target)
				a.PoC.Run()
				res := a.PoC.Report()
				if res["exploitable"].(bool) {
					exploitableList = append(exploitableList, res["target"].(string))
				}
			}(aa, blockers)
		}
		for i := 0; i < cap(blockers); i++ {
			blockers <- struct{}{}
		}
	} else {
		a.PoC.Set(a.Target)
		a.PoC.Run()
	}
}
