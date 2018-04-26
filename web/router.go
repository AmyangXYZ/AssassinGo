package web

import (
	"../config"
	"github.com/AmyangXYZ/sweetygo"
	"github.com/AmyangXYZ/sweetygo/middlewares"
)

var (
	requireJWTMap = map[string]string{
		"/api/*": "ALL",
		"/ws/*":  "ALL",
	}
)

// SetMiddlewares sets middlewares.
func SetMiddlewares(app *sweetygo.SweetyGo) {
	app.USE(middlewares.JWT("Header", config.SecretKey, requireJWTMap))
}

// SetRouter sets router.
func SetRouter(app *sweetygo.SweetyGo) {
	app.GET("/", index)
	app.POST("/token", signin)

	app.GET("/static/*files", static)

	app.POST("/api/target", setTarget)

	app.GET("/api/info/basic", basicInfo)
	app.GET("/api/info/whois", whois)
	app.GET("/api/info/cms", cmsDetect)
	app.GET("/api/info/honeypot", honeypot)

	app.GET("/ws/info/tracert", tracert)
	app.GET("/ws/info/port", portScan)
	app.GET("/ws/info/dirb", dirBrute)

	app.GET("/ws/attack/crawl", crawl)

	app.GET("/ws/attack/sqli", checkSQLi)
	app.GET("/ws/attack/xss", checkXSS)
	app.GET("/ws/attack/intrude", intrude)

	app.GET("/ws/seek", seek)

	app.GET("/api/poc", getPoCList)

	app.GET("/api/poc/:poc", runPoC)
	app.GET("/ws/poc/:poc", runSiblingPoC)
}
