package web

import (
	"github.com/AmyangXYZ/sweetygo"
)

// SetRouter sets router.
func SetRouter(app *sweetygo.SweetyGo) {
	app.GET("/", index)

	app.GET("/static/*files", static)

	app.POST("/api/assassin", newAssassin)
	app.POST("/api/assassin-dad", newAssassinDad)

	app.GET("/api/info/basic", basicInfo)
	app.GET("/api/info/whois", whois)
	app.GET("/api/info/cms", cmsDetect)
	app.GET("/api/info/honeypot", honeypot)

	app.GET("/ws/info/tracert", tracert)
	app.GET("/ws/info/port", portScan)
	app.GET("/ws/info/dirb", dirBrute)

	app.GET("/ws/info/crawl", crawl)

	app.GET("/ws/vul/sqli", checkSQLi)
	app.GET("/ws/vul/xss", checkXSS)
	app.GET("/ws/intrude", intrude)

	app.GET("/ws/seek", seek)

	app.GET("/api/poc", getPoCList)

	app.GET("/api/poc/:poc", runPoC)
	app.GET("/ws/poc/run", runPoCforSA)
}
