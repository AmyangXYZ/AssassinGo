package web

import (
	"../assassin"
	"github.com/AmyangXYZ/sweetygo"
)

var a *assassin.Assassin

// Run Web GUI.
func Run() {
	a = assassin.New("www.xxx.com")
	app := sweetygo.New("/home/amyang/Projects/AssassinGo/web", nil)
	SetRouter(app)
	app.RunServer(":8080")
}
