package config

import (
	"fmt"
	"os"
)

var (
	// RootDir of your app
	RootDir = "./web"

	// SecretKey computes sg_token
	SecretKey string

	// DB addr and passwd.
	DB string
)

func init() {
	DB = fmt.Sprintf("%v:%v@tcp(mariadb:3306)/%v?charset=utf8",
		os.Getenv("DB_User"),
		os.Getenv("DB_Passwd"),
		os.Getenv("DB_Db"))
	SecretKey = os.Getenv("SecretKey")
}
