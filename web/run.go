package web

import (
	"../assassin"
	"github.com/AmyangXYZ/sweetygo"
)

var a *assassin.Assassin
var dad *assassin.Dad

func init() {
	a = assassin.New()
	dad = assassin.NewDad()
}

// Run Web GUI.
func Run() {
	app := sweetygo.New("./web", nil)
	SetRouter(app)
	app.RunServer(":8000")
}
