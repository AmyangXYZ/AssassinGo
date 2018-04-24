package web

import (
	"net"
	"net/http"

	"../assassin"
	"../logger"
	"github.com/AmyangXYZ/sweetygo"
	"github.com/gorilla/websocket"
)

func index(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	ctx.Render(200, "index")
}

func static(ctx *sweetygo.Context) {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./web/static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
}

func newAssassin(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	a = assassin.New() // reset
	target := ctx.Param("target")
	a.SetTarget(target)
	ctx.JSON(200, 1, "success", nil)
}

func newAssassinDad(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	dad = assassin.NewDad() // reset
	targets := ctx.Param("targets")
	dad.SetTargets(targets)
	ctx.JSON(200, 1, "success", nil)
}

func basicInfo(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	a.Gatherers["basicInfo"].Set(a.Target)
	a.Gatherers["basicInfo"].Run()
	ret := a.Gatherers["basicInfo"].Report()
	ctx.JSON(200, 1, "success", ret)
}

func cmsDetect(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	a.Gatherers["cms"].Set(a.Target)
	a.Gatherers["cms"].Run()
	ret := a.Gatherers["cms"].Report()
	ctx.JSON(200, 1, "success", ret)
}

func whois(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	if net.ParseIP(a.Target).String() == a.Target {
		ctx.JSON(200, 0, "ip do not need whois", nil)
		return
	}
	a.Gatherers["whois"].Set(a.Target)
	a.Gatherers["whois"].Run()
	ret := a.Gatherers["whois"].Report()
	ctx.JSON(200, 1, "success", ret)
}

func honeypot(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	a.Gatherers["honeypot"].Set(a.Target)
	a.Gatherers["honeypot"].Run()
	ret := a.Gatherers["honeypot"].Report()
	ctx.JSON(200, 1, "success", ret)
}

func tracert(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["tracert"].Set(conn, a.Target)
	a.Gatherers["tracert"].Run()
	conn.Close()
}

func portScan(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
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
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := dirbMsg{}
	conn.ReadJSON(&m)
	a.Gatherers["dirb"].Set(conn, a.Target, m.GortCount)
	a.Gatherers["dirb"].Run()
	conn.Close()
}

func crawl(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["crawler"].Set(conn, a.Target, 4)
	a.Attackers["crawler"].Run()
	a.FuzzableURLs = a.Attackers["crawler"].Report()["fuzzableURLs"].([]string)
	conn.Close()
}

func checkSQLi(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["sqli"].Set(conn, a.FuzzableURLs)
	a.Attackers["sqli"].Run()
	conn.Close()
}

func checkXSS(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
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
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := intruderMsg{}
	err := conn.ReadJSON(&m)
	if err != nil {
		logger.Red.Println(err)
		conn.Close()
		return
	}
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
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := seekerMsg{}
	err := conn.ReadJSON(&m)
	if err != nil {
		logger.Red.Println(err)
		conn.Close()
		return
	}
	a.Seeker.Set(conn, m.Query, m.SE, m.MaxPage)
	a.Seeker.Run()
	conn.Close()
}

func getPoCList(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	pocList := []string{}
	for p := range a.PoC {
		pocList = append(pocList, p)
	}
	ret := map[string]interface{}{
		"poc_list": pocList,
	}
	ctx.JSON(200, 1, "success", ret)
}

func runPoC(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	pocName := ctx.Param("poc")
	if a.PoC[pocName] == nil {
		logger.Red.Println("No Such PoC")
		ctx.JSON(200, 0, "no such poc", nil)
		return
	}
	a.PoC[pocName].Set(a.Target)
	a.PoC[pocName].Run()
	ret := a.PoC[pocName].Report()
	ctx.JSON(200, 1, "success", ret)
}

type pocMsg struct {
	GortCount int `json:"gort_count"`
}

func runDadPoC(ctx *sweetygo.Context) {
	ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	pocName := ctx.Param("poc")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	if a.PoC[pocName] == nil {
		conn.WriteJSON(map[string]string{"message": "no such poc"})
		logger.Red.Println("No Such PoC")
		conn.Close()
		return
	}
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
	conn.Close()
}
