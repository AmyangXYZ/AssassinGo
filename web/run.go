package web

import (
	"../assassin"
	"github.com/AmyangXYZ/sweetygo"
)

var (
	daddy *assassin.Daddy
)

func init() {
	daddy = assassin.NewDaddy()
}

// Run Web GUI.
func Run() {
	app := sweetygo.New("./web", nil)
	SetMiddlewares(app)
	SetRouter(app)
	app.RunServer(":8000")
}
