package web

import (
	"../assassin"
	"../config"
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
	app := sweetygo.New(config.RootDir, nil)
	SetMiddlewares(app)
	SetRouter(app)
	app.RunServer(":8000")
}
