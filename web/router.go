package web

import (
	"github.com/AmyangXYZ/sweetygo"
)

// SetRouter sets router.
func SetRouter(app *sweetygo.SweetyGo) {
	app.GET("/", index)

	app.GET("/static/*files", static)

	app.POST("/api/target", newAssassin)

	app.GET("/ws/info/basic", basicInfo)
	app.GET("/ws/info/whois", whois)
	app.GET("/ws/info/tracert", tracert)
	app.GET("/wx/info/cms", cmsDetect)
	app.GET("/wx/info/port", portScan)

	app.GET("/ws/info/crawl", crawl)

	app.GET("/ws/vul/sqli", checkSQLi)
	app.GET("/ws/vul/xss", checkXSS)
	app.GET("/ws/intrude", intrude)

	app.GET("/ws/seek", seek)

	app.GET("/api/poc", getPOCs)
	app.POST("/api/poc", setPOC)
	app.GET("/ws/poc/run", runPOC)
}
