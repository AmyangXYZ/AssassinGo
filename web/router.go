package web

import (
	"../config"
	"github.com/AmyangXYZ/sgo"
	"github.com/AmyangXYZ/sgo/middlewares"
)

// SetMiddlewares sets middlewares.
func SetMiddlewares(app *sgo.sgo) {
	// cors
	app.USE(func(ctx *sgo.Context) error {
		ctx.Resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		ctx.Resp.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Next()
		return nil
	})

	jwtSkipper := func(ctx *sgo.Context) bool {
		if (len(ctx.Path()) > 5 && ctx.Path()[0:5] == "/api/") ||
			(len(ctx.Path()) > 4 && ctx.Path()[0:4] == "/ws/") {
			return false
		}
		return true
	}
	app.USE(middlewares.JWT("Cookie", config.SecretKey, jwtSkipper))

}

// SetRouter sets router.
func SetRouter(app *sgo.sgo) {
	app.GET("/", index)
	app.POST("/token", signin)

	app.GET("/static/*files", static)

	app.POST("/api/target", setTarget)

	app.GET("/api/info/basic", basicInfo)
	app.GET("/api/info/bypasscf", bypassCF)
	app.GET("/api/info/whois", whois)
	app.GET("/api/info/cms", cmsDetect)
	app.GET("/api/info/honeypot", honeypot)

	app.GET("/ws/info/tracert", tracert)
	app.GET("/ws/info/port", portScan)
	app.GET("/ws/info/subdomain", subDomainScan)
	app.GET("/ws/info/dirb", dirBrute)

	app.GET("/ws/attack/crawl", crawl)

	app.GET("/ws/attack/sqli", checkSQLi)
	app.GET("/ws/attack/xss", checkXSS)
	app.GET("/ws/attack/intrude", intrude)

	app.GET("/ws/attack/ssh", sshBrute)

	app.GET("/ws/seek", seek)

	app.GET("/api/poc", getPoCList)

	app.GET("/api/poc/:poc", runPoC)
	app.GET("/ws/poc/:poc", runSiblingPoC)
}
