package web

import (
	"assassingo/assassin"

	"github.com/AmyangXYZ/sgo"
)

var (
	daddy *assassin.Daddy
)

func init() {
	daddy = assassin.NewDaddy()
}

// Run Web GUI.
func Run() {
	app := sgo.New()
	app.SetTemplates("/web/templates", nil)
	SetMiddlewares(app)
	SetRouter(app)
	app.Run(":8000")
}
