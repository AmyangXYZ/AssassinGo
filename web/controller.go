package web

import (
	"net"
	"net/http"
	"time"

	"../assassin"
	"../config"
	"../logger"
	"../poc"
	"github.com/AmyangXYZ/sweetygo"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

func index(ctx *sweetygo.Context) error {
	return ctx.Render(200, "index")
}

func static(ctx *sweetygo.Context) error {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./web/static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
	return nil
}

func signin(ctx *sweetygo.Context) error {
	if ctx.Param("username") != "" && ctx.Param("password") != "" {
		username := ctx.Param("username")
		password := getPassword(username)

		if password == ctx.Param("password") {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["username"] = username
			claims["exp"] = time.Now().Add(time.Hour * 36).Unix()
			t, _ := token.SignedString([]byte(config.SecretKey))
			ctx.SetCookie("SG_Token", t)
			ctx.JSON(200, 1, "success", map[string]string{"SG_Token": t})

			a := assassin.New()
			s := assassin.NewSiblings()
			daddy.Son[username] = a
			daddy.Sibling[username] = s
			logger.Green.Println(username, "Has Signed In")

			return nil
		}
		return ctx.JSON(200, 0, "Username or Password Error.", nil)
	}
	return nil
}

func setTarget(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	if target := ctx.Param("target"); target != "" {
		daddy.Son[usr].SetTarget(target)
		return ctx.JSON(200, 1, "success", nil)
	}
	if targets := ctx.Param("targets"); targets != "" {
		daddy.Sibling[usr].SetTargets(targets)
		return ctx.JSON(200, 1, "success", nil)
	}
	return nil
}

func basicInfo(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	a.Gatherers["basicInfo"].Set(a.Target)
	a.Gatherers["basicInfo"].Run()
	ret := a.Gatherers["basicInfo"].Report()
	return ctx.JSON(200, 1, "success", ret)
}

func bypassCF(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	a.Gatherers["bypassCF"].Set(a.Target)
	a.Gatherers["bypassCF"].Run()
	ret := a.Gatherers["bypassCF"].Report()
	return ctx.JSON(200, 1, "success", ret)
}

func cmsDetect(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	a.Gatherers["cms"].Set(a.Target)
	a.Gatherers["cms"].Run()
	ret := a.Gatherers["cms"].Report()
	return ctx.JSON(200, 1, "success", ret)
}

func whois(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	if net.ParseIP(a.Target).String() == a.Target {
		return ctx.JSON(200, 0, "ip do not need whois", nil)
	}
	a.Gatherers["whois"].Set(a.Target)
	a.Gatherers["whois"].Run()
	ret := a.Gatherers["whois"].Report()
	return ctx.JSON(200, 1, "success", ret)
}

func honeypot(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	a.Gatherers["honeypot"].Set(a.Target)
	a.Gatherers["honeypot"].Run()
	ret := a.Gatherers["honeypot"].Report()
	return ctx.JSON(200, 1, "success", ret)
}

func tracert(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["tracert"].Set(conn, a.Target)
	a.Gatherers["tracert"].Run()
	conn.Close()
	return nil
}

func portScan(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["port"].Set(conn, a.Target)
	a.Gatherers["port"].Run()
	conn.Close()
	return nil
}

func subDomainScan(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["subdomain"].Set(conn, a.Target)
	a.Gatherers["subdomain"].Run()
	conn.Close()
	return nil
}

type dirbMsg struct {
	// Payload   string `json:"payload"`
	Concurrency int `json:"concurrency"`
}

func dirBrute(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := dirbMsg{}
	conn.ReadJSON(&m)
	a.Gatherers["dirb"].Set(conn, a.Target, m.Concurrency)
	a.Gatherers["dirb"].Run()
	conn.Close()
	return nil
}

func crawl(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["crawler"].Set(conn, a.Target, 4)
	a.Attackers["crawler"].Run()
	a.FuzzableURLs = a.Attackers["crawler"].Report()["fuzzableURLs"].([]string)
	conn.Close()
	return nil
}

func checkSQLi(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["sqli"].Set(conn, a.FuzzableURLs)
	a.Attackers["sqli"].Run()
	conn.Close()
	return nil
}

func checkXSS(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["xss"].Set(conn, a.FuzzableURLs)
	a.Attackers["xss"].Run()
	conn.Close()
	return nil
}

type intruderMsg struct {
	Header      string `json:"header"`
	Payload     string `json:"payload"`
	Concurrency int    `json:"concurrency"`
}

func intrude(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := intruderMsg{}
	err := conn.ReadJSON(&m)
	if err != nil {
		logger.Red.Println(err)
		conn.Close()
		return err
	}
	a.Attackers["intruder"].Set(conn, a.Target, m.Header, m.Payload, m.Concurrency)
	a.Attackers["intruder"].Run()
	conn.Close()
	return nil
}

type sshMsg struct {
	Port        string `json:"port"`
	Concurrency int    `json:"concurrency"`
}

func sshBrute(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := sshMsg{}
	err := conn.ReadJSON(&m)
	if err != nil {
		logger.Red.Println(err)
		conn.Close()
		return err
	}
	a.Attackers["ssh"].Set(conn, a.Target, m.Port, m.Concurrency)
	a.Attackers["ssh"].Run()
	conn.Close()
	return nil
}

type seekerMsg struct {
	Query   string `json:"query"`
	SE      string `json:"se"`
	MaxPage int    `json:"max_page"`
}

func seek(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := seekerMsg{}
	err := conn.ReadJSON(&m)
	if err != nil {
		logger.Red.Println(err)
		conn.Close()
		return err
	}
	a.Seeker.Set(conn, m.Query, m.SE, m.MaxPage)
	a.Seeker.Run()
	conn.Close()
	return nil
}

func getPoCList(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	pocList := map[string]poc.Intro{}
	for k, v := range a.PoC {
		pocList[k] = v.Info()
	}
	ret := map[string]interface{}{
		"poc_list": pocList,
	}
	return ctx.JSON(200, 1, "success", ret)
}

func runPoC(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	a := daddy.Son[usr]
	pocName := ctx.Param("poc")
	if a.PoC[pocName] == nil {
		logger.Red.Println("No Such PoC")
		return ctx.JSON(200, 0, "no such poc", nil)
	}
	a.PoC[pocName].Set(a.Target)
	a.PoC[pocName].Run()
	ret := a.PoC[pocName].Report()
	return ctx.JSON(200, 1, "success", ret)
}

type pocMsg struct {
	Concurrency int `json:"concurrency"`
}

func runSiblingPoC(ctx *sweetygo.Context) error {
	usr := ctx.Get("userInfo").(*jwt.Token).Claims.(jwt.MapClaims)["username"].(string)
	sibling := daddy.Sibling[usr]
	pocName := ctx.Param("poc")
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	if daddy.Son[usr].PoC[pocName] == nil {
		conn.WriteJSON(map[string]string{"message": "no such poc"})
		logger.Red.Println("No Such PoC")
		conn.Close()
		return nil
	}
	sibling.MuxConn.Conn = conn
	m := pocMsg{}
	conn.ReadJSON(&m)

	blockers := make(chan struct{}, m.Concurrency)
	for _, a := range sibling.Siblings {
		blockers <- struct{}{}
		go func(a *assassin.Assassin, blocker chan struct{}) {
			defer func() { <-blocker }()
			a.PoC[pocName].Set(a.Target)
			a.PoC[pocName].Run()
			ret := a.PoC[pocName].Report()
			if ret["exploitable"].(bool) {
				sibling.MuxConn.Send(ret)
				sibling.ExploitableHosts = append(sibling.ExploitableHosts, ret["host"].(string))
			}
		}(a, blockers)
	}

	for i := 0; i < cap(blockers); i++ {
		blockers <- struct{}{}
	}
	conn.Close()
	return nil
}
