package web

import (
	"github.com/AmyangXYZ/sweetygo"
)

// SetRouter sets router.
func SetRouter(app *sweetygo.SweetyGo) {
	app.GET("/", home)

	app.POST("/api/target", newAssassin)
	app.GET("/api/info/basic", basicInfo)
	app.GET("/api/info/cms", cmsDetect)
	app.GET("/api/info/port", portScan)
	app.GET("/api/info/goohack", nil)

	app.GET("/api/crawl", crawl)

	app.GET("/api/vul/sqli", nil)
	app.GET("/api/vul/xss", nil)

	app.GET("/api/poc/xxx", nil)
}
