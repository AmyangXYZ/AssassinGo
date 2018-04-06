package web

import (
	"github.com/AmyangXYZ/sweetygo"
)

// SetRouter sets router.
func SetRouter(app *sweetygo.SweetyGo) {
	app.GET("/", index)

	app.GET("/static/*files", static)

	app.POST("/api/target", newAssassin)
	app.POST("/api/targets", setTargets)

	app.GET("/api/info/basic", basicInfo)
	app.GET("/api/info/cms", cmsDetect)
	app.GET("/api/info/port", portScan)
	app.GET("/api/info/goohack", nil)

	app.GET("/api/ws/crawl", wsCrawl)
	app.GET("/api/ws/intrude", intrude)

	app.GET("/api/vul/sqli", checkSQLi)
	app.GET("/api/vul/xss", checkXSS)

	app.GET("/api/poc", getPOCs)
	app.POST("/api/poc", setPOC)
	app.GET("/api/poc/run", runPOC)

}
