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
	app.GET("/ws/info/goohack", nil)

	app.GET("/ws/crawl", crawl)
	app.GET("/ws/intrude", intrude)

	app.GET("/ws/vul/sqli", checkSQLi)
	app.GET("/ws/vul/xss", checkXSS)

	app.GET("/ws/seek", seek)

	app.GET("/api/poc", getPOCs)
	app.POST("/api/poc", setPOC)
	app.GET("/ws/poc/run", runPOC)

}
