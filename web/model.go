package web

import (
	"database/sql"
	"fmt"
	"time"

	//
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "ag:agpasswd@tcp(mariadb:3306)/ag?charset=utf8")

	// mysql image starts need time.
	for {
		err := db.Ping()
		if err == nil {
			fmt.Println("db ok!")
			break
		}
		fmt.Println(err)
		time.Sleep(2 * time.Second)
	}

	// https://github.com/go-sql-driver/mysql/issues/674
	db.SetMaxIdleConns(0)

	db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT(10) NOT NULL AUTO_INCREMENT,
			username VARCHAR(16) NULL DEFAULT NULL,
			password VARCHAR(64) NULL DEFAULT NULL,
			email VARCHAR(64) NULL DEFAULT NULL,
			PRIMARY KEY (id)
		);`)
}

func getPassword(username string) string {
	var password string
	rows, _ := db.Query("SELECT password from users where username=?", username)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&password)
	}
	return password
}
