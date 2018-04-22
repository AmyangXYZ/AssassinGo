package web

import (
	"net/http"

	"../assassin"
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
	a = assassin.New() // reset
	target := ctx.Param("target")
	a.SetTarget(target)
	ctx.JSON(201, nil, "success")
}

func newAssassinDad(ctx *sweetygo.Context) {
	dad = assassin.NewDad() // reset
	targets := ctx.Param("targets")
	dad.SetTargets(targets)
	ctx.JSON(201, nil, "success")
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
	a.Seeker.Set(conn, m.Query, m.SE, m.MaxPage)
	a.Seeker.Run()
	conn.Close()
}

func getPoCList(ctx *sweetygo.Context) {
	pocList := []string{}
	for p := range a.PoC {
		pocList = append(pocList, p)
	}
	ret := map[string]interface{}{
		"poc_list": pocList,
	}
	ctx.JSON(200, ret, "success")
}

func runPoC(ctx *sweetygo.Context) {
	pocName := ctx.Param("poc")
	a.PoC[pocName].Set(a.Target)
	a.PoC[pocName].Run()
	ret := a.PoC[pocName].Report()
	ctx.JSON(200, ret, "success")
}

type pocMsg struct {
	GortCount int `json:"gort_count"`
}

func runDadPoC(ctx *sweetygo.Context) {
	pocName := ctx.Param("poc")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	dad.MuxConn.Conn = conn
	m := pocMsg{}
	conn.ReadJSON(&m)

	blockers := make(chan struct{}, m.GortCount)
	for _, son := range dad.Sons {
		blockers <- struct{}{}
		go func(son *assassin.Assassin, blocker chan struct{}) {
			defer func() { <-blocker }()
			son.PoC[pocName].Set(son.Target)
			son.PoC[pocName].Run()
			ret := son.PoC[pocName].Report()
			if ret["exploitable"].(bool) {
				dad.MuxConn.Send(ret)
				dad.ExploitableHosts = append(dad.ExploitableHosts, ret["host"].(string))
			}
		}(son, blockers)
	}

	for i := 0; i < cap(blockers); i++ {
		blockers <- struct{}{}
	}
}
